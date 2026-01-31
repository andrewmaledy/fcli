package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// State represents the current phase of the interactive selection flow.
type State int

const (
	StateLoading  State = iota
	StateTable          // Navigate, select, confirm
	StateConfirm        // y/N confirmation
	StateDeleting       // Spinner while deleting
	StateResults        // Show results
	StateError          // Show error
)

// fetchDoneMsg is sent when the FetchItemsFunc completes.
type fetchDoneMsg struct {
	items []ListItem
	err   error
}

// deleteDoneMsg is sent when the DeleteItemsFunc completes.
type deleteDoneMsg struct {
	results []string
}

// SelectTableModel is the Bubble Tea model for the enhanced interactive selection.
type SelectTableModel struct {
	state   State
	title   string
	columns []ColumnDef

	// Data
	items   []ListItem
	fetchFn FetchItemsFunc
	deleteFn DeleteItemsFunc

	// Cursor and expansion
	cursor         int
	visibleRows    []visibleRow
	expandedParent int // index into items[], or -1
	scrollOffset   int

	// Selection: keys are "i" for parent, "i.c" for child
	selections map[string]bool

	// Terminal dimensions
	width  int
	height int

	// UI components
	spinner spinner.Model

	// Output
	results  []string
	errMsg   string
	quitting bool

	// Public: populated on exit for callers to inspect
	SelectedSelections []Selection
}

// NewSelectTableModel creates a new interactive selection model.
func NewSelectTableModel(title string, columns []ColumnDef, fetchFn FetchItemsFunc, deleteFn DeleteItemsFunc) SelectTableModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(ColorCyan)

	return SelectTableModel{
		state:          StateLoading,
		title:          title,
		columns:        columns,
		spinner:        s,
		selections:     make(map[string]bool),
		fetchFn:        fetchFn,
		deleteFn:       deleteFn,
		expandedParent: -1,
	}
}

func (m SelectTableModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			items, err := m.fetchFn()
			return fetchDoneMsg{items: items, err: err}
		},
	)
}

func (m SelectTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case fetchDoneMsg:
		if msg.err != nil {
			m.state = StateError
			m.errMsg = msg.err.Error()
			return m, nil
		}
		if len(msg.items) == 0 {
			m.state = StateError
			m.errMsg = "No items found."
			return m, nil
		}

		m.items = msg.items
		m.cursor = 0
		m.expandedParent = -1
		m.visibleRows = m.buildVisibleRows()
		m.state = StateTable
		return m, nil

	case deleteDoneMsg:
		m.state = StateResults
		m.results = msg.results
		return m, nil

	case spinner.TickMsg:
		if m.state == StateLoading || m.state == StateDeleting {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case tea.KeyMsg:
		switch m.state {
		case StateTable:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					m.ensureCursorVisible()
				}
				return m, nil
			case "down", "j":
				if m.cursor < len(m.visibleRows)-1 {
					m.cursor++
					m.ensureCursorVisible()
				}
				return m, nil
			case "right", "l":
				m.expandCurrent()
				m.ensureCursorVisible()
				return m, nil
			case "left", "h":
				m.collapseCurrent()
				m.ensureCursorVisible()
				return m, nil
			case " ":
				m.toggleSelection()
				return m, nil
			case "enter":
				sels := m.buildSelections()
				if len(sels) > 0 {
					if m.deleteFn != nil {
						m.state = StateConfirm
					} else {
						m.SelectedSelections = sels
						m.quitting = true
						return m, tea.Quit
					}
				}
				return m, nil
			case "q", "esc", "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}

		case StateConfirm:
			switch msg.String() {
			case "y", "Y":
				m.state = StateDeleting
				sels := m.buildSelections()
				m.SelectedSelections = sels
				deleteFn := m.deleteFn
				return m, tea.Batch(
					m.spinner.Tick,
					func() tea.Msg {
						results := deleteFn(sels)
						return deleteDoneMsg{results: results}
					},
				)
			case "n", "N", "esc":
				m.state = StateTable
				return m, nil
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}

		case StateResults, StateError:
			switch msg.String() {
			case "q", "enter", "esc", "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m SelectTableModel) View() string {
	if m.quitting {
		return ""
	}

	switch m.state {
	case StateLoading:
		return fmt.Sprintf("\n  %s Loading %s...\n\n", m.spinner.View(), m.title)
	case StateTable:
		return m.renderTableView()
	case StateConfirm:
		return m.renderConfirmView()
	case StateDeleting:
		return fmt.Sprintf("\n  %s Deleting...\n\n", m.spinner.View())
	case StateResults:
		return m.renderResultsView()
	case StateError:
		return fmt.Sprintf("\n%s\n\n  %s\n\n  %s\n",
			ErrorStyle.Render("Error"),
			m.errMsg,
			DimStyle.Render("Press [enter] or [q] to exit"),
		)
	}
	return ""
}

// --- Visible Row Management ---

func (m *SelectTableModel) buildVisibleRows() []visibleRow {
	var rows []visibleRow
	for i, item := range m.items {
		rows = append(rows, visibleRow{itemIdx: i, childIdx: -1})
		if i == m.expandedParent && len(item.Children) > 0 {
			for c := range item.Children {
				rows = append(rows, visibleRow{itemIdx: i, childIdx: c})
			}
		}
	}
	return rows
}

// expandCurrent expands the parent under the cursor and moves cursor to its first child.
func (m *SelectTableModel) expandCurrent() {
	if m.cursor < 0 || m.cursor >= len(m.visibleRows) {
		return
	}
	row := m.visibleRows[m.cursor]

	// Only expand if cursor is on a parent with children
	if row.childIdx != -1 || len(m.items[row.itemIdx].Children) == 0 {
		return
	}

	m.expandedParent = row.itemIdx
	m.visibleRows = m.buildVisibleRows()

	// Move cursor to first child (it's right after the parent in visibleRows)
	for i, vr := range m.visibleRows {
		if vr.itemIdx == row.itemIdx && vr.childIdx == 0 {
			m.cursor = i
			return
		}
	}
}

// collapseCurrent collapses the expanded parent and returns cursor to it.
func (m *SelectTableModel) collapseCurrent() {
	if m.cursor < 0 || m.cursor >= len(m.visibleRows) {
		return
	}
	row := m.visibleRows[m.cursor]

	// If on a child, collapse back to its parent
	if row.childIdx >= 0 {
		m.expandedParent = -1
		m.visibleRows = m.buildVisibleRows()
		// Find the parent row
		for i, vr := range m.visibleRows {
			if vr.itemIdx == row.itemIdx && vr.childIdx == -1 {
				m.cursor = i
				return
			}
		}
		return
	}

	// If on an expanded parent, just collapse it
	if row.itemIdx == m.expandedParent {
		m.expandedParent = -1
		m.visibleRows = m.buildVisibleRows()
		// Cursor stays on same parent, just find its new position
		for i, vr := range m.visibleRows {
			if vr.itemIdx == row.itemIdx && vr.childIdx == -1 {
				m.cursor = i
				return
			}
		}
	}
}

func (m *SelectTableModel) ensureCursorVisible() {
	viewportHeight := m.viewportHeight()
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
	if m.cursor >= m.scrollOffset+viewportHeight {
		m.scrollOffset = m.cursor - viewportHeight + 1
	}
}

func (m SelectTableModel) viewportHeight() int {
	h := m.height - 8 // header, footer, borders
	if h < 5 {
		h = 20
	}
	return h
}

// --- Selection Logic ---

func parentKey(i int) string   { return fmt.Sprintf("%d", i) }
func childKey(i, c int) string { return fmt.Sprintf("%d.%d", i, c) }

func (m *SelectTableModel) toggleSelection() {
	if m.cursor < 0 || m.cursor >= len(m.visibleRows) {
		return
	}
	row := m.visibleRows[m.cursor]

	if row.childIdx == -1 {
		// Toggling a parent
		key := parentKey(row.itemIdx)
		if m.selections[key] {
			delete(m.selections, key)
			for c := range m.items[row.itemIdx].Children {
				delete(m.selections, childKey(row.itemIdx, c))
			}
		} else {
			m.selections[key] = true
			for c := range m.items[row.itemIdx].Children {
				delete(m.selections, childKey(row.itemIdx, c))
			}
		}
	} else {
		// Toggling a child
		cKey := childKey(row.itemIdx, row.childIdx)
		pKey := parentKey(row.itemIdx)

		if m.selections[pKey] {
			// Parent selected -> deselecting this child means demote
			delete(m.selections, pKey)
			for c := range m.items[row.itemIdx].Children {
				if c != row.childIdx {
					m.selections[childKey(row.itemIdx, c)] = true
				}
			}
		} else if m.selections[cKey] {
			delete(m.selections, cKey)
		} else {
			m.selections[cKey] = true
			// Auto-promote: if all children now selected, promote to parent
			allSelected := true
			for c := range m.items[row.itemIdx].Children {
				if !m.selections[childKey(row.itemIdx, c)] {
					allSelected = false
					break
				}
			}
			if allSelected && len(m.items[row.itemIdx].Children) > 0 {
				m.selections[pKey] = true
				for c := range m.items[row.itemIdx].Children {
					delete(m.selections, childKey(row.itemIdx, c))
				}
			}
		}
	}
}

func (m SelectTableModel) isSelected(itemIdx, childIdx int) bool {
	if m.selections[parentKey(itemIdx)] {
		return true
	}
	if childIdx >= 0 {
		return m.selections[childKey(itemIdx, childIdx)]
	}
	return false
}

func (m SelectTableModel) buildSelections() []Selection {
	var sels []Selection
	for i, item := range m.items {
		if m.selections[parentKey(i)] {
			sels = append(sels, Selection{ItemIndex: i, ChildIndex: -1})
			continue
		}
		for c := range item.Children {
			if m.selections[childKey(i, c)] {
				sels = append(sels, Selection{ItemIndex: i, ChildIndex: c})
			}
		}
	}
	return sels
}

// --- Rendering ---

func (m SelectTableModel) renderTableView() string {
	totalWidth := m.width
	if totalWidth < 80 {
		totalWidth = 120
	}
	totalHeight := m.height
	if totalHeight < 10 {
		totalHeight = 30
	}

	listWidth := totalWidth*55/100 - 2
	detailWidth := totalWidth - listWidth - 5

	contentHeight := totalHeight - 6

	leftPanel := m.renderListPanel(listWidth, contentHeight)
	rightPanel := m.renderDetailPanel(detailWidth, contentHeight)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	header := TitleStyle.Render(m.title)
	selCount := len(m.buildSelections())
	selInfo := ""
	if selCount > 0 {
		selInfo = "  " + WarningStyle.Render(fmt.Sprintf("%d selected", selCount))
	}
	help := DimStyle.Render("  [space] toggle  [\u2192] expand  [\u2190] collapse  [enter] confirm  [q] quit")

	return fmt.Sprintf("\n%s%s\n\n%s\n\n%s\n", header, selInfo, content, help)
}

func (m SelectTableModel) renderListPanel(width, height int) string {
	var lines []string

	// Header row
	lines = append(lines, m.renderHeaderRow(width-4))

	// Viewport
	vpHeight := height - 3
	startRow := m.scrollOffset
	endRow := startRow + vpHeight
	if endRow > len(m.visibleRows) {
		endRow = len(m.visibleRows)
	}

	for i := startRow; i < endRow; i++ {
		vr := m.visibleRows[i]
		isCursor := (i == m.cursor)
		selected := m.isSelected(vr.itemIdx, vr.childIdx)
		lines = append(lines, m.renderRow(vr, selected, isCursor, width-4))
	}

	// Pad remaining lines
	for len(lines) < vpHeight+1 {
		lines = append(lines, "")
	}

	joined := strings.Join(lines, "\n")

	style := ListBorderStyle.Copy().Width(width).Height(height)
	return style.Render(joined)
}

func (m SelectTableModel) renderHeaderRow(width int) string {
	cols := m.formatColumnsHeader(width - 8)
	line := fmt.Sprintf("       %s", cols)
	return SubtitleStyle.Bold(true).Render(line)
}

func (m SelectTableModel) renderRow(vr visibleRow, selected, cursor bool, width int) string {
	item := m.items[vr.itemIdx]

	checkbox := "[ ]"
	if selected {
		checkbox = SuccessStyle.Render("[x]")
	}

	if vr.childIdx >= 0 {
		// Child row
		child := item.Children[vr.childIdx]
		isLast := vr.childIdx == len(item.Children)-1

		connector := TreeConnectorStyle.Render("  \u251C ")
		if isLast {
			connector = TreeConnectorStyle.Render("  \u2514 ")
		}

		// Compact child display: name, episodes, size
		name := child.Columns[0]
		episodes := ""
		size := ""
		if len(child.Columns) > 3 {
			episodes = child.Columns[3]
		}
		if len(child.Columns) > 6 {
			size = child.Columns[6]
		}

		childContent := fmt.Sprintf("%-18s %8s  %8s", name, episodes, size)
		line := fmt.Sprintf("%s%s %s", connector, checkbox, childContent)

		if cursor {
			return CursorStyle.Copy().Width(width).Render(line)
		}
		return line
	}

	// Parent row
	indicator := " "
	if len(item.Children) > 0 {
		if vr.itemIdx == m.expandedParent {
			indicator = DimStyle.Render("\u25BC")
		} else {
			indicator = DimStyle.Render("\u25B6")
		}
	}

	cols := m.formatColumns(item.Columns, width-8)
	line := fmt.Sprintf("%s %s %s", checkbox, indicator, cols)

	if cursor {
		return CursorStyle.Copy().Width(width).Render(line)
	}
	return line
}

func (m SelectTableModel) formatColumnsHeader(availableWidth int) string {
	totalDef := 0
	for _, c := range m.columns {
		totalDef += c.Width
	}
	if totalDef == 0 {
		return ""
	}

	var parts []string
	for _, col := range m.columns {
		colWidth := col.Width * availableWidth / totalDef
		if colWidth < 1 {
			colWidth = 1
		}
		title := col.Title
		if len(title) > colWidth {
			title = title[:colWidth]
		}
		parts = append(parts, fmt.Sprintf("%-*s", colWidth, title))
	}
	return strings.Join(parts, "")
}

func (m SelectTableModel) formatColumns(cols []string, availableWidth int) string {
	totalDef := 0
	for _, c := range m.columns {
		totalDef += c.Width
	}
	if totalDef == 0 {
		return ""
	}

	var parts []string
	for i, col := range m.columns {
		colWidth := col.Width * availableWidth / totalDef
		if colWidth < 1 {
			colWidth = 1
		}
		val := ""
		if i < len(cols) {
			val = cols[i]
		}
		if len(val) > colWidth {
			val = val[:colWidth-1] + "\u2026"
		}
		parts = append(parts, fmt.Sprintf("%-*s", colWidth, val))
	}
	return strings.Join(parts, "")
}

// --- Detail Panel ---

func (m SelectTableModel) renderDetailPanel(width, height int) string {
	detailSection := m.renderDetailSection(width - 4)
	selectionSection := m.renderSelectionSummary(width - 4)

	content := detailSection
	if selectionSection != "" {
		content += "\n\n" + selectionSection
	}

	style := ListBorderStyle.Copy().Width(width).Height(height).Padding(0, 1)
	return style.Render(content)
}

func (m SelectTableModel) renderDetailSection(width int) string {
	if m.cursor < 0 || m.cursor >= len(m.visibleRows) {
		return ""
	}

	vr := m.visibleRows[m.cursor]
	item := m.items[vr.itemIdx]

	header := TitleStyle.Render("\u2500\u2500 Details \u2500\u2500")

	var detail []KeyValue
	if vr.childIdx >= 0 && vr.childIdx < len(item.Children) {
		detail = item.Children[vr.childIdx].Detail
	} else {
		detail = item.Detail
	}

	var lines []string
	lines = append(lines, header)
	lines = append(lines, "")

	maxKeyWidth := 0
	for _, kv := range detail {
		if len(kv.Key) > maxKeyWidth {
			maxKeyWidth = len(kv.Key)
		}
	}

	for _, kv := range detail {
		key := DetailKeyStyle.Render(fmt.Sprintf("%-*s", maxKeyWidth+1, kv.Key+":"))
		val := DetailValueStyle.Render(kv.Value)
		lines = append(lines, fmt.Sprintf(" %s %s", key, val))
	}

	return strings.Join(lines, "\n")
}

func (m SelectTableModel) renderSelectionSummary(width int) string {
	selections := m.buildSelections()
	if len(selections) == 0 {
		return ""
	}

	header := SelectionHeaderStyle.Render(fmt.Sprintf("\u2500\u2500 Selected (%d) \u2500\u2500", len(selections)))

	var lines []string
	lines = append(lines, header)
	lines = append(lines, "")

	for _, sel := range selections {
		item := m.items[sel.ItemIndex]
		label := ""
		if sel.ChildIndex == -1 {
			label = item.Label
			if len(item.Children) > 0 {
				label += " (Full)"
			}
		} else if sel.ChildIndex < len(item.Children) {
			label = item.Children[sel.ChildIndex].Label
		}

		bullet := SelectionBulletStyle.Render("\u2022")
		lines = append(lines, fmt.Sprintf(" %s %s", bullet, label))
	}

	return strings.Join(lines, "\n")
}

// --- Confirm / Results Views ---

func (m SelectTableModel) renderConfirmView() string {
	sels := m.buildSelections()
	header := WarningStyle.Render(fmt.Sprintf("Delete %d selected items? (y/N)", len(sels)))
	var items string
	for _, sel := range sels {
		item := m.items[sel.ItemIndex]
		label := item.Label
		if sel.ChildIndex >= 0 && sel.ChildIndex < len(item.Children) {
			label = fmt.Sprintf("%s - %s", item.Columns[0], item.Children[sel.ChildIndex].Label)
		} else if len(item.Children) > 0 {
			label += " (Full series)"
		}
		items += fmt.Sprintf("  - %s\n", label)
	}
	return fmt.Sprintf("\n%s\n\n%s\n", header, items)
}

func (m SelectTableModel) renderResultsView() string {
	header := TitleStyle.Render("Results")
	var body string
	for _, result := range m.results {
		body += fmt.Sprintf("  %s\n", result)
	}
	help := DimStyle.Render("  Press [enter] or [q] to exit")
	return fmt.Sprintf("\n%s\n\n%s\n%s\n", header, body, help)
}

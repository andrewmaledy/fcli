package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// RenderStaticTable creates a styled table string from columns and rows.
// This does NOT launch a Bubble Tea program -- it renders a string for printing.
func RenderStaticTable(columns []table.Column, rows []table.Row) string {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(len(rows)+1),
		table.WithFocused(false),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ColorGray).
		BorderBottom(true).
		Bold(true).
		Foreground(ColorCyan)
	s.Cell = s.Cell.Padding(0, 1)
	// No visual selection highlight for static display
	s.Selected = s.Selected.
		Foreground(lipgloss.NoColor{}).
		Background(lipgloss.NoColor{}).
		Bold(false)

	t.SetStyles(s)

	return TableBorderStyle.Render(t.View())
}

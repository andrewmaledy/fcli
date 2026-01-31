package tui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	ColorRed     = lipgloss.Color("#FF5555")
	ColorGreen   = lipgloss.Color("#50FA7B")
	ColorYellow  = lipgloss.Color("#F1FA8C")
	ColorBlue    = lipgloss.Color("#6272A4")
	ColorCyan    = lipgloss.Color("#8BE9FD")
	ColorMagenta = lipgloss.Color("#FF79C6")
	ColorWhite   = lipgloss.Color("#F8F8F2")
	ColorGray    = lipgloss.Color("#6272A4")
	ColorDimGray = lipgloss.Color("#44475A")
)

// Text styles
var (
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(ColorCyan)
	SubtitleStyle = lipgloss.NewStyle().Foreground(ColorGray)
	SuccessStyle  = lipgloss.NewStyle().Foreground(ColorGreen)
	ErrorStyle    = lipgloss.NewStyle().Foreground(ColorRed)
	WarningStyle  = lipgloss.NewStyle().Foreground(ColorYellow)
	DimStyle      = lipgloss.NewStyle().Foreground(ColorDimGray)
)

// Table styles
var (
	TableBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(ColorGray)
)

// Detail panel styles
var (
	DetailKeyStyle   = lipgloss.NewStyle().Foreground(ColorCyan).Bold(true)
	DetailValueStyle = lipgloss.NewStyle().Foreground(ColorWhite)
)

// Selection summary styles
var (
	SelectionHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorYellow)
	SelectionBulletStyle = lipgloss.NewStyle().Foreground(ColorYellow)
)

// Tree and cursor styles
var (
	TreeConnectorStyle = lipgloss.NewStyle().Foreground(ColorDimGray)
	CursorStyle        = lipgloss.NewStyle().
				Foreground(ColorWhite).
				Background(lipgloss.Color("#44475A")).
				Bold(true)
	ListBorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorGray)
)

// Error box styles
var (
	ErrorBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(ColorRed).
			Padding(1, 2).
			Width(65)
	ErrorBoxTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorRed).
				Align(lipgloss.Center)
)

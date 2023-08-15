package show

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WindowSize stores the size of the terminal
var WindowSize tea.WindowSizeMsg

// STYLES

// Viewport: The viewport of the tui (my:2, mx:40)
var Viewport = lipgloss.NewStyle().Margin(1, 2)

// Container: container styling (width: 80, py: 0, px: 5)
var Container = lipgloss.NewStyle().Width(80).Padding(0, 5)

// WidgetContainer: container for tables, and graphs (py:2, px:2)
var WidgetContainer = lipgloss.NewStyle().Padding(2, 2)

// SquareBorder: Style to draw a border around a section
var SquareBorder = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(Color)

// TextContainer: container for text
var TextContainer = lipgloss.NewStyle().Padding(1, 1)

// TableTitle: The style for table titles (width:25, align-horizontal:center, bold:true)
var TableTitle = lipgloss.NewStyle().Width(25).AlignHorizontal(lipgloss.Center).Bold(true)

// Color: the color palette (Light: #000000, Dark: #FF4500)
var Color = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FF4500"}

// ActiveStyle: table when selected (border:normal, border-foreground:#FF4500)
var ActiveStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#FF4500"))

// InactiveStyle: table when not selected (border:normal, border-foreground:#FFFFFF)
var InactiveStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#FFFFFF"))

package styles

import "github.com/charmbracelet/lipgloss"

var Styles = InitStyles()

var (
	TabsHeight         = 3
	TableHeaderHeight  = 2
	FooterHeight       = 1
	StatusBarHeight    = 1
	ExpandedHelpHeight = 8
	CellPadding        = 1

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}
	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}
)

type styles struct {
	Colors struct {
	}
	Common struct {
		Footer lipgloss.Style
	}
	Tabs struct {
		Tab       lipgloss.Style
		ActiveTab lipgloss.Style
		TabsRow   lipgloss.Style
	}
	StatusBar struct {
		Bar          lipgloss.Style
		Subscription lipgloss.Style
		User         lipgloss.Style
		Text         lipgloss.Style
	}
	Section struct {
		Container lipgloss.Style
		Spinner   lipgloss.Style
	}
	Table struct {
		Table struct {
			Header   lipgloss.Style
			Cell     lipgloss.Style
			Selected lipgloss.Style
		}
	}
	Spinner struct {
		Spinner lipgloss.Style
	}
	PopUp struct {
		Button         lipgloss.Style
		Text           lipgloss.Style
		Box            lipgloss.Style
		SelectedButton lipgloss.Style
	}
}

func InitStyles() styles {
	var styles styles

	styles.Common.Footer = lipgloss.NewStyle().
		Height(FooterHeight)

	styles.Tabs.Tab = lipgloss.NewStyle().
		Faint(true).
		Padding(0, 1).
		Border(tabBorder)
	styles.Tabs.ActiveTab = styles.Tabs.Tab.Copy().
		Faint(false).
		Bold(true).
		Border(activeTabBorder)
	styles.Tabs.TabsRow = lipgloss.NewStyle().
		Height(TabsHeight)

	styles.StatusBar.Bar = lipgloss.NewStyle().
		Height(StatusBarHeight).
		Background(Theme.BackgroundSecondary)
	styles.StatusBar.Subscription = lipgloss.NewStyle().
		Inherit(styles.StatusBar.Bar).
		Padding(0, 1).
		Bold(true).
		Background(Theme.PrimaryAccent).
		Foreground(Theme.PrimaryText)
	styles.StatusBar.User = lipgloss.NewStyle().
		Inherit(styles.StatusBar.Bar).
		Padding(0, 1)
	styles.StatusBar.Text = lipgloss.NewStyle().
		Inherit(styles.StatusBar.Bar)

	styles.Section.Container = lipgloss.NewStyle().
		Padding(0, 0)

	styles.Table.Table.Selected = lipgloss.NewStyle().Bold(true).Background(Theme.PrimaryAccent)
	styles.Table.Table.Header = lipgloss.NewStyle().
		Bold(false).
		Padding(0, 1).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("240")).
		BorderStyle(lipgloss.NormalBorder())
	styles.Table.Table.Cell = lipgloss.NewStyle().Padding(0, CellPadding)

	styles.Spinner.Spinner = lipgloss.NewStyle().Foreground(Theme.PrimaryAccent)

	styles.PopUp.Button = lipgloss.NewStyle().
		Foreground(Theme.PrimaryText).
		Background(lipgloss.Color("#888B7E")).
		Padding(0, 3).
		MarginRight(2).
		MarginTop(1)

	styles.PopUp.SelectedButton = lipgloss.NewStyle().
		Foreground(Theme.PrimaryText).
		Background(Theme.PrimaryAccent).
		Padding(0, 3).
		MarginTop(1).
		MarginRight(2).
		Underline(true)
	styles.PopUp.Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Theme.PrimaryAccent).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	return styles
}

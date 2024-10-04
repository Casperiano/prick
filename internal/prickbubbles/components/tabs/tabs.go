package tabs

import (
	"prick/internal/prickbubbles/cmds"
	"prick/internal/prickbubbles/context"
	"prick/internal/prickbubbles/keys"
	"prick/internal/prickbubbles/styles"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var tabs = []string{"󰮄 Resource Groups", "󰆧 Resources", " Pricks"}

type Model struct {
	ctx       *context.BubbleContext
	id        string
	ActiveTab int
	TabCount  int
}

func New(ctx *context.BubbleContext) Model {
	return Model{
		ctx:       ctx,
		id:        zone.NewPrefix(),
		ActiveTab: 0,
		TabCount:  3,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Choose):
			if m.ActiveTab == m.TabCount-1 {
				break
			}
			m.ActiveTab++
			cmd = append(cmd, func() tea.Msg { return cmds.TabSwitchedMsg{TabToLoad: m.ActiveTab} })
		case key.Matches(msg, keys.Keys.Quit):
			m.ActiveTab--
			cmd = append(cmd, func() tea.Msg { return cmds.TabSwitchedMsg{TabsToRefresh: []int{m.ActiveTab + 1}} })
		}
	case tea.MouseMsg:
		if msg.Action != tea.MouseAction(tea.MouseButtonLeft) {
			break
		}

		for i, tab := range tabs {
			if zone.Get(m.id + tab).InBounds(msg) {
				if i < m.ActiveTab {
					var tabsToRefresh []int
					for j := m.ActiveTab; j > i; j-- {
						tabsToRefresh = append(tabsToRefresh, j)
					}
					cmd = append(cmd, func() tea.Msg { return cmds.TabSwitchedMsg{TabsToRefresh: tabsToRefresh} })
				}

				if i > m.ActiveTab {
					if i-m.ActiveTab > 1 {
						break
					}
					cmd = append(cmd, func() tea.Msg { return cmds.TabSwitchedMsg{TabToLoad: i} })
				}

				m.ActiveTab = i
				break
			}
		}
	}
	return m, tea.Batch(cmd...)
}

func (m Model) View() string {
	var renderedTabs []string
	for i, tab := range tabs {
		if m.ActiveTab == i {
			renderedTabs = append(renderedTabs,
				zone.Mark(m.id+tab, styles.Styles.Tabs.ActiveTab.Render(tab)),
			)
		} else {
			renderedTabs = append(renderedTabs,
				zone.Mark(m.id+tab, styles.Styles.Tabs.Tab.Render(tab)),
			)
		}
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderedTabs...,
	)

	widthFill := styles.Styles.Tabs.Tab.Copy().BorderLeft(false).BorderRight(false).BorderTop(false).
		Render(strings.Repeat(" ", max(0, m.ctx.ScreenWidth-lipgloss.Width(row)-2)))

	return styles.Styles.Tabs.TabsRow.Copy().
		Render(lipgloss.JoinHorizontal(lipgloss.Bottom, row, widthFill))
}

func (m *Model) SetActiveTab(id int) {
	m.ActiveTab = id
}

func (m *Model) UpdateContext(ctx *context.BubbleContext) {
	m.ctx = ctx
}

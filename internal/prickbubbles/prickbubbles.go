package prickbubbles

import (
	"log"
	prick "prick/internal/prick"
	"prick/internal/prick/config"
	command "prick/internal/prickbubbles/cmds"
	"prick/internal/prickbubbles/components/footer"
	"prick/internal/prickbubbles/components/pricksection"
	"prick/internal/prickbubbles/components/rgsection"
	"prick/internal/prickbubbles/components/rsection"
	"prick/internal/prickbubbles/components/statusbar"
	"prick/internal/prickbubbles/components/tabs"
	"prick/internal/prickbubbles/context"
	"prick/internal/prickbubbles/keys"
	"prick/internal/prickbubbles/styles"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

type Model struct {
	footer       footer.Model
	statusbar    statusbar.Model
	keys         keys.KeyMap
	tabs         tabs.Model
	rgsection    rgsection.Model
	rsection     rsection.Model
	pricksection pricksection.Model
	ctx          context.BubbleContext
}

func New(api *prick.Api, config *config.Config) Model {
	zone.NewGlobal()
	ctx := context.BubbleContext{Config: config, Api: api}

	tabs := tabs.New(&ctx)
	rgsection := rgsection.New(&ctx)
	rsection := rsection.New(&ctx)
	pricksection := pricksection.New(&ctx)
	footer := footer.New(&ctx)
	statusbar := statusbar.New(&ctx)

	return Model{
		ctx:          ctx,
		keys:         keys.Keys,
		tabs:         tabs,
		rgsection:    rgsection,
		rsection:     rsection,
		pricksection: pricksection,
		footer:       footer,
		statusbar:    statusbar,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		command.FetchResourceGroups(m.ctx.Api),
		command.InitConfig,
		m.statusbar.Init(),
		m.rgsection.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	cmds := m.updateShell(msg)

	switch msg := msg.(type) {
	case command.InitMsg:
		cmds = append(cmds, command.Refresh(5))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.updateContentHeight()
		case key.Matches(msg, m.keys.Quit):
			// Only quit if not searching and on the first tab
			if !m.rgsection.IsSearching && m.tabs.ActiveTab < 0 {
				cmd = tea.Quit
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		m.updateWindowSize(msg)
	case command.TabSwitchedMsg:
		if msg.TabToLoad > 0 {
			switch msg.TabToLoad {
			case 1:
				cmds = append(cmds, m.rsection.Init())
			case 2:
				cmds = append(cmds, m.pricksection.Init())
			}
		}
		for _, tab := range msg.TabsToRefresh {
			switch tab {
			case 1:
				m.rsection.Clear()
			case 2:
				m.pricksection.Clear()
			}
		}

	case command.TickEventMsg:
		// Here we call APIs to refresh data

		if m.rgsection.Table.Cursor() >= 0 {
			selectedResourceGroup, err := m.rgsection.GetSelectedResourceGroup()
			if err != nil {
				log.Default().Printf("Error getting selected resource group: %v", err)
			}
			cmds = append(cmds,
				command.FetchResources(m.ctx.Api, selectedResourceGroup),
			)
		}

		if m.rsection.Table.Cursor() >= 0 {
			selectedResource := m.rsection.GetSelectedResource()
			if selectedResource != nil {
				cmds = append(cmds,
					command.FetchPokes(m.ctx.Api, selectedResource),
				)
			}
		}

		cmds = append(
			cmds, command.Refresh(2),
		)
	}

	switch m.tabs.ActiveTab {
	case 0:
		m.rgsection.UpdateContext(&m.ctx)
		m.rgsection, cmd = m.rgsection.Update(msg)
	case 1:
		m.rsection.UpdateContext(&m.ctx)
		m.rsection, cmd = m.rsection.Update(msg)
	case 2:
		m.pricksection.UpdateContext(&m.ctx)
		m.pricksection, cmd = m.pricksection.Update(msg)
	}

	return m, tea.Batch(append(cmds, cmd)...)
}

func (m Model) View() string {
	view := strings.Builder{}

	view.WriteString(m.tabs.View() + "\n")

	switch m.tabs.ActiveTab {
	case 0:
		view.WriteString(m.rgsection.View())
	case 1:
		view.WriteString(m.rsection.View())
	case 2:
		view.WriteString(m.pricksection.View())
	}

	view.WriteString("\n")
	view.WriteString(m.footer.View() + "\n")
	view.WriteString(m.statusbar.View())

	return zone.Scan(view.String())
}

func (m *Model) updateShell(msg tea.Msg) []tea.Cmd {
	var (
		cmdTabs      tea.Cmd
		cmdFooter    tea.Cmd
		cmdStatusbar tea.Cmd
	)

	m.tabs.UpdateContext(&m.ctx)
	m.footer.UpdateContext(&m.ctx)
	m.statusbar.UpdateContext(&m.ctx)

	// If we're searching, we don't want to switch tabs using Esc and Enter
	if !m.IsSearching() {
		m.tabs, cmdTabs = m.tabs.Update(msg)
	}
	m.footer, cmdFooter = m.footer.Update(msg)
	m.statusbar, cmdStatusbar = m.statusbar.Update(msg)

	return []tea.Cmd{
		cmdTabs,
		cmdFooter,
		cmdStatusbar,
	}
}

func (m *Model) IsSearching() bool {
	return m.rgsection.IsSearching || m.rsection.IsSearching
}

func (m *Model) updateWindowSize(msg tea.WindowSizeMsg) {
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
	m.updateContentHeight()
}

func (m *Model) updateContentHeight() {
	contentHeight := m.ctx.ScreenHeight - styles.TabsHeight - styles.TableHeaderHeight - styles.FooterHeight - styles.StatusBarHeight

	if m.footer.Help.ShowAll {
		contentHeight -= styles.ExpandedHelpHeight
	}
	if m.rgsection.IsSearching {
		contentHeight -= 1
	}

	m.ctx.ContentHeight = contentHeight
}

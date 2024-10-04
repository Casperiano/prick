package rgsection

import (
	"fmt"
	"prick/internal/prickbubbles/cmds"
	"prick/internal/prickbubbles/components/search"
	"prick/internal/prickbubbles/context"
	"prick/internal/prickbubbles/keys"
	"prick/internal/prickbubbles/styles"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ctx         *context.BubbleContext
	Id          int
	Table       table.Model
	Spinner     spinner.Model
	Rows        []table.Row
	SearchBar   search.Model
	IsSearching bool
	IsLoading   bool
}

func New(ctx *context.BubbleContext) Model {
	rgtable := table.New(
		table.WithColumns(getColumns(ctx.ScreenWidth, false)),
		table.WithFocused(true),
		table.WithWidth(ctx.ScreenWidth),
		table.WithStyles(styles.Styles.Table.Table),
	)

	spinner := spinner.New(
		spinner.WithSpinner(spinner.Points),
		spinner.WithStyle(styles.Styles.Spinner.Spinner),
	)

	search := search.New(ctx)

	return Model{
		ctx:         ctx,
		Spinner:     spinner,
		Table:       rgtable,
		SearchBar:   search,
		IsSearching: false,
		IsLoading:   true,
	}
}

func getColumns(screenWidth int, isSearching bool) []table.Column {
	nameCol := "Name"
	if isSearching {
		nameCol = nameCol + " "
	}

	regionWidth := 20
	return []table.Column{
		{Title: nameCol, Width: screenWidth - regionWidth - styles.CellPadding*2},
		{Title: "Region", Width: regionWidth - styles.CellPadding*2},
	}
}

func (m Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.IsSearching {
			switch {
			case key.Matches(msg, keys.Keys.Help):
				return m, cmd // Shortcut so that question marks don't end up in the search field
			case key.Matches(msg, keys.Keys.Quit):
				cmd = m.SetIsSearching(false)
			case key.Matches(msg, keys.Keys.Choose):
				cmd = m.SetIsSearching(false)
			}
		} else {
			switch {
			case key.Matches(msg, keys.Keys.Search):
				cmd = m.SetIsSearching(true)
				return m, cmd
			case key.Matches(msg, keys.Keys.Up):
				m.Table.MoveUp(1)
			case key.Matches(msg, keys.Keys.Down):
				m.Table.MoveDown(1)
			}
		}

	case cmds.ResourceGroupsFetchedMsg:
		m.updateRows(msg)
		m.IsLoading = false
	case spinner.TickMsg:
		if m.IsLoading {
			m.Spinner, cmd = m.Spinner.Update(msg)
			return m, cmd
		}
	}

	if m.IsSearching {
		m.SearchBar, cmd = m.SearchBar.Update(msg)
		m.Table.GotoTop()
	}
	m.filterRows()
	if m.SearchBar.TextInput.Value() != "" {
		m.Table.SetColumns(getColumns(m.ctx.ScreenWidth, true))
	}

	return m, cmd
}

func (m Model) View() string {
	view := strings.Builder{}

	if m.IsLoading {
		view.WriteString(lipgloss.Place(
			m.ctx.ScreenWidth,
			m.ctx.ContentHeight+2,
			lipgloss.Center,
			lipgloss.Center,
			m.Spinner.View(),
		))
	} else {
		view.WriteString(m.Table.View())
	}

	if m.IsSearching {
		view.WriteString("\n" + m.SearchBar.View())
	}

	return view.String()
}

func (m *Model) updateRows(msg cmds.ResourceGroupsFetchedMsg) {
	var rows []table.Row
	// slices.Sort(msg.ResourceGroupNames)
	for _, rg := range msg.ResourceGroups {
		rows = append(rows, table.Row{*rg.Name, *rg.Location})
	}
	m.Rows = rows
}

func (m *Model) filterRows() {
	var rows []table.Row
	searchValue := m.SearchBar.TextInput.Value()

	for _, row := range m.Rows {
		if strings.Contains(strings.ToLower(row[0]), strings.ToLower(searchValue)) {
			rows = append(rows, row)
		}
	}

	m.Table.SetRows(rows)
}

func (m *Model) UpdateContext(ctx *context.BubbleContext) {
	m.ctx = ctx
	m.Table.SetHeight(m.ctx.ContentHeight)
	m.Table.SetWidth(m.ctx.ScreenWidth)
	m.Table.SetColumns(getColumns(m.ctx.ScreenWidth, false))
}

func (m *Model) SetIsSearching(val bool) tea.Cmd {
	m.IsSearching = val
	if val {
		m.SearchBar.SetFocus(true)
		m.ctx.ContentHeight -= 1
		m.Table.SetHeight(m.ctx.ContentHeight)
		return m.SearchBar.Init()
	} else {
		m.SearchBar.SetFocus(false)
		m.ctx.ContentHeight += 1
		m.Table.SetHeight(m.ctx.ContentHeight)
		return nil
	}
}

func (m *Model) GetSelectedResourceGroup() (string, error) {
	if m.Table.Cursor() < 0 || len(m.Table.Rows()) == 0 {
		return "", fmt.Errorf("no resource group selected")
	}

	return m.Table.Rows()[m.Table.Cursor()][0], nil

}

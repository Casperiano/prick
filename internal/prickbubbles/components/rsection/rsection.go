package rsection

import (
	"prick/internal/prick/interfaces"
	"prick/internal/prickbubbles/components/search"
	"prick/internal/prickbubbles/context"
	"prick/internal/prickbubbles/keys"
	"prick/internal/prickbubbles/styles"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	command "prick/internal/prickbubbles/cmds"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ctx         *context.BubbleContext
	Spinner     spinner.Model
	Table       table.Model
	Rows        []table.Row
	Prickables  []interfaces.Prickable
	SearchBar   search.Model
	IsSearching bool
	IsLoading   bool
	Lock        *sync.Mutex
}

func getColumns(screenWidth int, isSearching bool) []table.Column {
	searchColName := "Name"
	if isSearching {
		searchColName = searchColName + " "
	}

	return []table.Column{
		{Title: searchColName, Width: (screenWidth-20)/2 - styles.CellPadding*2},
		{Title: "Type", Width: (screenWidth-20)/2 - styles.CellPadding*2},
		{Title: "Region", Width: 20 - styles.CellPadding*2},
	}
}

func New(ctx *context.BubbleContext) Model {

	rtable := table.New(
		table.WithColumns(getColumns(ctx.ScreenWidth, false)),
		table.WithStyles(styles.Styles.Table.Table),
		table.WithWidth(ctx.ScreenWidth),
		table.WithFocused(true),
	)
	search := search.New(ctx)

	spinner := spinner.New(
		spinner.WithSpinner(spinner.Points),
		spinner.WithStyle(styles.Styles.Spinner.Spinner),
	)

	return Model{
		ctx:         ctx,
		Spinner:     spinner,
		Table:       rtable,
		SearchBar:   search,
		IsSearching: false,
		IsLoading:   true,
		Lock:        &sync.Mutex{},
	}
}

func (m Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case command.ResourcesFetchedMsg:
		m.Prickables = msg.Prickables
		m.updateRows(msg)
		m.IsLoading = false
	case spinner.TickMsg:
		if m.IsLoading {
			m.Spinner, cmd = m.Spinner.Update(msg)
			return m, cmd
		}
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

	}
	m.ctx.SelectedResource = m.GetSelectedResource()

	if m.IsSearching {
		m.SearchBar, cmd = m.SearchBar.Update(msg)
		m.Table.GotoTop()
	}
	m.filterRows()
	// If the search state contains a value, we want to indicate that there is an active filter with a 
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

func (m *Model) filterRows() {
	m.Lock.Lock()
	defer m.Lock.Unlock()

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

func (m *Model) GetSelectedResource() interfaces.Prickable {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if m.Table.Cursor() < 0 || len(m.Table.Rows()) == 0 {
		return nil
	}

	// Keep getting index out of range sometimes, even with the lock, use % hack.
	currRow := m.Table.Rows()[m.Table.Cursor()%len(m.Table.Rows())]
	name, rtype, region := currRow[0], currRow[1], currRow[2]

	for _, p := range m.Prickables {
		if p.GetName() == name && p.GetLocation() == region && string(p.GetType()) == rtype {
			return p
		}
	}

	return m.Prickables[m.Table.Cursor()]
}

func (m *Model) updateRows(msg command.ResourcesFetchedMsg) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	var rows []table.Row
	for _, p := range msg.Prickables {
		rows = append(rows, table.Row{p.GetName(), string(p.GetType()), p.GetLocation()})
	}
	m.Rows = rows
}

func (m *Model) Clear() {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.Rows = []table.Row{}
	m.Table.SetRows(m.Rows)
	m.Prickables = []interfaces.Prickable{}
	m.IsLoading = true
}

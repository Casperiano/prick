package pricksection

import (
	"fmt"
	"prick/internal/prick/common"
	"prick/internal/prickbubbles/components/search"
	"prick/internal/prickbubbles/context"
	"prick/internal/prickbubbles/keys"
	"prick/internal/prickbubbles/styles"
	"strings"

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
	Pokes       []*common.Poke
	Rows        []table.Row
	SearchBar   search.Model
	IsSearching bool
	IsLoading   bool
	PopUp       PopUp
}

type PopUp struct {
	IsActive bool
	Confirm  bool
	Action   Action
}

type Action int

const (
	Add Action = iota
	Patch
)

func getColumns(screenWidth int, isSearching bool) []table.Column {
	searchColName := "Name"
	if isSearching {
		searchColName = searchColName + " "
	}
	return []table.Column{
		{Title: searchColName, Width: 40 - styles.CellPadding*2},
		{Title: "StartIp", Width: (screenWidth-40)/2 - styles.CellPadding*2},
		{Title: "EndIp", Width: (screenWidth-40)/2 - styles.CellPadding*2},
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
		PopUp:       PopUp{IsActive: false, Confirm: false, Action: Add},
	}
}

func (m Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case error:
		return m, tea.Quit
	case command.PokesFetchedMsg:
		m.Pokes = msg.Pokes
		m.updateRows(msg)
		m.IsLoading = false
	case command.PokedMsg:
		m.IsLoading = false
	case command.PatchedMsg:
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
		} else if m.PopUp.IsActive {
			switch {
			case key.Matches(msg, keys.Keys.NextSection) || key.Matches(msg, keys.Keys.PrevSection):
				m.PopUp.Confirm = !m.PopUp.Confirm
			case key.Matches(msg, keys.Keys.Quit):
				m.PopUp.IsActive = false
			case key.Matches(msg, keys.Keys.Choose):
				m.PopUp.IsActive = false
				if m.PopUp.Confirm {
					m.PopUp.Confirm = false
					m.IsLoading = true
					switch {
					case m.PopUp.Action == Add:
						return m, tea.Batch(command.Poke(m.ctx.Api, m.ctx.SelectedResource), m.Spinner.Tick)
					case m.PopUp.Action == Patch:
						return m, tea.Batch(command.Patch(m.ctx.Api, m.ctx.SelectedResource), m.Spinner.Tick)
					}
				}
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
			case key.Matches(msg, keys.Keys.Add):
				m.PopUp.IsActive = true
				m.PopUp.Action = Add
			case key.Matches(msg, keys.Keys.Patch):
				m.PopUp.IsActive = true
				m.PopUp.Action = Patch
			}

		}

	}

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

	if m.PopUp.IsActive {
		var okButton, cancelButton string
		if m.PopUp.Confirm {
			okButton = styles.Styles.PopUp.SelectedButton.Render("Yes")
			cancelButton = styles.Styles.PopUp.Button.Render("No")
		} else {
			okButton = styles.Styles.PopUp.Button.Render("Yes")
			cancelButton = styles.Styles.PopUp.SelectedButton.Render("No")
		}

		var popUpMessage string
		if m.PopUp.Action == Add {
			popUpMessage = "add"
		} else {
			popUpMessage = "patch"
		}
		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(fmt.Sprintf("Are you sure you want to %s your ip?", popUpMessage))
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)

		ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)
		dialog := styles.Styles.PopUp.Box.Render(ui)

		return PlaceOverlay((m.ctx.ScreenWidth-50)/2, m.ctx.ContentHeight/2, dialog, view.String(), false)

	}

	return view.String()
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

func (m *Model) updateRows(msg command.PokesFetchedMsg) {
	var rows []table.Row
	for _, poke := range msg.Pokes {
		rows = append(rows, table.Row{poke.Name, poke.StartIpAddress, poke.EndIpAddress})
	}
	m.Rows = rows
}

func (m *Model) Clear() {
	m.Rows = []table.Row{}
	m.Table.SetRows(m.Rows)
	m.Pokes = []*common.Poke{}
	m.IsLoading = true
}

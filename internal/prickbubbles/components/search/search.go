package search

import (
	"prick/internal/prickbubbles/context"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ctx       *context.BubbleContext
	TextInput textinput.Model
}

func New(ctx *context.BubbleContext) Model {
	textInput := textinput.New()
	textInput.Prompt = "  "
	textInput.Placeholder = "Search..."
	textInput.Cursor.BlinkSpeed = time.Millisecond * 500
	textInput.Cursor.Blink = true

	return Model{TextInput: textInput}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	m.TextInput, cmd = m.TextInput.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.TextInput.View()
}

func (m *Model) UpdateContext(ctx *context.BubbleContext) {
	m.ctx = ctx
}

func (m *Model) SetFocus(val bool) {
	if val {
		m.TextInput.TextStyle = m.TextInput.TextStyle.Copy().Faint(false)
		m.TextInput.CursorEnd()
		m.TextInput.Focus()
	} else {
		m.TextInput.TextStyle = m.TextInput.TextStyle.Copy().Faint(true)
		m.TextInput.CursorStart()
		m.TextInput.Blur()
	}
}

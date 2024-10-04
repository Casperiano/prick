package footer

import (
	"prick/internal/prickbubbles/context"
	"prick/internal/prickbubbles/keys"
	"prick/internal/prickbubbles/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ctx  *context.BubbleContext
	Help help.Model
}

func New(ctx *context.BubbleContext) Model {
	help := help.New()
	help.ShowAll = false

	return Model{
		ctx:  ctx,
		Help: help,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Help):
			m.Help.ShowAll = !m.Help.ShowAll
		}
	}

	return m, nil
}

func (m Model) View() string {
	footer := styles.Styles.Common.Footer.Copy().
		Render(lipgloss.JoinVertical(lipgloss.Top, m.Help.View(keys.Keys)))

	return footer
}

func (m *Model) UpdateContext(ctx *context.BubbleContext) {
	m.ctx = ctx
}

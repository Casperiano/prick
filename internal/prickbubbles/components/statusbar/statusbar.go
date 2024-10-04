package statusbar

import (
	"prick/internal/prick/common"
	"prick/internal/prickbubbles/cmds"
	"prick/internal/prickbubbles/context"
	"prick/internal/prickbubbles/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ctx          *context.BubbleContext
	subscription string
	user         string
	tasks        string
}

func New(ctx *context.BubbleContext) Model {
	return Model{
		ctx: ctx,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		accountInfo, err := common.GetAzAccountInfo()
		if err != nil {
			return err
		}
		return cmds.AccountInfoFetchedMsg{AccountInfo: accountInfo}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case cmds.AccountInfoFetchedMsg:
		m.subscription = msg.AccountInfo.SubscriptionName
		m.user = msg.AccountInfo.User.Name
	}
	return m, nil
}

func (m Model) View() string {
	subscription := styles.Styles.StatusBar.Subscription.Copy().Render(" " + m.subscription)
	user := styles.Styles.StatusBar.User.Copy().Render(" " + m.user)
	help := styles.Styles.StatusBar.Text.Copy().Padding(0, 1).Render("? help")

	spacer := styles.Styles.StatusBar.Text.Copy().
		Width(m.ctx.ScreenWidth - lipgloss.Width(subscription) - lipgloss.Width(user) - lipgloss.Width(m.tasks) - lipgloss.Width(help)).
		Render("")

	statusbar := styles.Styles.StatusBar.Bar.Copy().
		Width(m.ctx.ScreenWidth).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Top,
			subscription,
			user,
			spacer,
			m.tasks,
			help,
		))

	return statusbar
}

func (m *Model) SetTasks(tasks string) {
	m.tasks = tasks
}

func (m *Model) UpdateContext(ctx *context.BubbleContext) {
	m.ctx = ctx
}

package cmds

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickEventMsg time.Time

func Refresh(refreshInterval int) tea.Cmd {
	return tea.Tick(
		time.Second*time.Duration(refreshInterval),
		func(t time.Time) tea.Msg {
			return TickEventMsg(t)
		},
	)
}

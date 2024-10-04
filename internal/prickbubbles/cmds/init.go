package cmds

import tea "github.com/charmbracelet/bubbletea"

type Config struct {
}

type InitMsg struct {
	Config Config
}

func InitConfig() tea.Msg {
	// Here config will be parsed
	return InitMsg{Config: Config{}}
}

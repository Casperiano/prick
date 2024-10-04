package cmds

import (
	"prick/internal/prick"
	"prick/internal/prick/common"
	"prick/internal/prick/interfaces"

	tea "github.com/charmbracelet/bubbletea"
)

type PokesFetchedMsg struct {
	Pokes []*common.Poke
}

func FetchPokes(api *prick.Api, r interfaces.Prickable) tea.Cmd {
	return func() tea.Msg {
		pokes, err := r.ListPokes(api)
		if err != nil {
			return err
		}

		return PokesFetchedMsg{Pokes: pokes}
	}
}

type PokedMsg struct{}

func Poke(api *prick.Api, r interfaces.Prickable) tea.Cmd {
	return func() tea.Msg {
		err := r.Poke(api)
		if err != nil {
			return err
		}
		return PokedMsg{}
	}
}

type PatchedMsg struct{}

func Patch(api *prick.Api, r interfaces.Prickable) tea.Cmd {
	return func() tea.Msg {
		err := r.Patch(api)
		if err != nil {
			return err
		}
		return PatchedMsg{}
	}
}

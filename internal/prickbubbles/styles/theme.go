package styles

import "github.com/charmbracelet/lipgloss"

type theme struct {
	PrimaryText         lipgloss.Color
	PrimaryAccent       lipgloss.Color
	Green               lipgloss.Color
	BackgroundSecondary lipgloss.Color
	BackgroundHighlight lipgloss.Color
}

var Theme = &theme{
	BackgroundSecondary: lipgloss.Color("236"),
	BackgroundHighlight: lipgloss.Color("8"),
	PrimaryText:         lipgloss.Color("7"),
	PrimaryAccent:       lipgloss.Color("4"),
	Green:               lipgloss.Color("10"),
}

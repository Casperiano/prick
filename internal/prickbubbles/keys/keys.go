package keys

import "github.com/charmbracelet/bubbles/key"

var Keys = KeyMap{
	Up:          key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:        key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Choose:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Search:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	Refresh:     key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
	PrevSection: key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "prev section")),
	NextSection: key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "next section")),
	Help:        key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	Quit:        key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", "quit")),
	Add:         key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
	Patch:       key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "patch")),
	PatchAll:    key.NewBinding(key.WithKeys("P"), key.WithHelp("P", "patch all")),
}

type KeyMap struct {
	Up          key.Binding
	Down        key.Binding
	Choose      key.Binding
	Search      key.Binding
	Refresh     key.Binding
	PrevSection key.Binding
	NextSection key.Binding
	Help        key.Binding
	Quit        key.Binding
	Add         key.Binding
	Patch       key.Binding
	PatchAll    key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		k.Up,
		k.Down,
		k.Choose,
		k.Search,
		k.Refresh,
		k.PrevSection,
		k.NextSection,
		k.Help,
		k.Quit,
		k.Add,
		k.Patch,
		k.PatchAll,
	}}
}

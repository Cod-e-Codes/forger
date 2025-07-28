package core

import tea "github.com/charmbracelet/bubbletea"

// Plugin is the interface every plugin must implement.
type Plugin interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Plugin, tea.Cmd)
	View() string
	Name() string
}

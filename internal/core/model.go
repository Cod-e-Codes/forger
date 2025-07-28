package core

import (
	"forger/internal/plugins/codesleuth"
	"forger/internal/plugins/ignoregrets"
	"forger/internal/plugins/marchat"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is the core Bubble Tea model for Forger.
type Model struct {
	Plugins    map[string]Plugin
	Active     string
	Overlay    Plugin
	Context    *Context
	LoadErrors []string
	Styles     lipgloss.Style
}

// NewModel constructs a Model with default styling and an empty Context.
func NewModel() Model {
	return Model{
		Context:    &Context{GlobalState: make(map[string]interface{})},
		Styles:     lipgloss.NewStyle().Padding(1).Border(lipgloss.NormalBorder()),
		LoadErrors: nil,
	}
}

func (m Model) Init() tea.Cmd {
	// Call Init() for all plugins, not just the active one
	var cmds []tea.Cmd
	for _, plugin := range m.Plugins {
		cmds = append(cmds, plugin.Init())
	}

	// Return a command that runs all the init commands
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Overlay routing
	if m.Overlay != nil {
		updated, cmd := m.Overlay.Update(msg)
		m.Overlay = updated
		if key, ok := msg.(tea.KeyMsg); ok && (key.String() == "c" || key.String() == "esc") {
			m.Overlay = nil
		}
		return m, cmd
	}

	// Check if this is an availability message that should be sent to all plugins
	switch msg.(type) {
	case marchat.ServerCheckMsg, codesleuth.AvailabilityMsg, ignoregrets.AvailabilityMsg:
		// Route availability messages to all plugins
		var cmds []tea.Cmd
		for name, plugin := range m.Plugins {
			updated, cmd := plugin.Update(msg)
			m.Plugins[name] = updated
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		if len(cmds) > 0 {
			return m, tea.Batch(cmds...)
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "c":
			if chat, ok := m.Plugins["marchat"]; ok {
				m.Overlay = chat
			}
			return m, nil
		case "esc":
			m.Overlay = nil
			return m, nil
		case "tab":
			// Use tab to switch between plugins instead of up/down
			m.Active = NextPluginKey(m.Plugins, m.Active)
			return m, nil
		case "shift+tab":
			// Use shift+tab to go backwards
			m.Active = PrevPluginKey(m.Plugins, m.Active)
			return m, nil
		}
	}

	// Update active plugin - let it handle all keys including up/down
	if p, ok := m.Plugins[m.Active]; ok {
		updated, cmd := p.Update(msg)
		m.Plugins[m.Active] = updated
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder

	// Show load errors if any
	if len(m.LoadErrors) > 0 {
		sb.WriteString("Errors loading plugins:\n")
		for _, e := range m.LoadErrors {
			sb.WriteString("  • " + e + "\n")
		}
		sb.WriteString("\n")
	}

	// Sidebar
	for _, name := range SortedPluginNames(m.Plugins) {
		prefix := "  "
		if name == m.Active {
			prefix = "> "
		}
		sb.WriteString(prefix + name + "\n")
	}
	sb.WriteString("\n")

	// Main content
	mainContent := ""
	if m.Overlay != nil {
		mainContent = m.Overlay.View()
	} else if p, ok := m.Plugins[m.Active]; ok {
		mainContent = p.View()
	} else {
		mainContent = "No active plugin."
	}
	styledContent := m.Styles.Render(mainContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, sb.String(), styledContent)
}

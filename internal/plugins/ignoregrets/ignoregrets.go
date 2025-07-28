package ignoregrets

import (
	"forger/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

type Plugin struct {
	ctx *types.Context
}

func New(ctx *types.Context) types.Plugin {
	return &Plugin{ctx: ctx}
}

func (p *Plugin) Init() tea.Cmd {
	return nil
}

func (p *Plugin) Update(msg tea.Msg) (types.Plugin, tea.Cmd) {
	return p, nil
}

func (p *Plugin) View() string {
	return `
┌─ IgnoreGrets ──────────────────────────────────────────────┐
│                                                             │
│  Snapshot Management                                        │
│                                                             │
│  • Create snapshot: Ctrl+S                                  │
│  • List snapshots: L                                       │
│  • Restore snapshot: R                                     │
│  • Delete snapshot: D                                       │
│                                                             │
│  Current snapshots:                                         │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ No snapshots available                                 │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                             │
│  Press 'h' for help, 'q' to quit                          │
└─────────────────────────────────────────────────────────────┘`
}

func (p *Plugin) Name() string {
	return "ignoregrets"
}

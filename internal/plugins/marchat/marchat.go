package marchat

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
┌─ MarChat ───────────────────────────────────────────────────┐
│                                                             │
│  Terminal Chat Interface                                    │
│                                                             │
│  • Send message: Enter                                     │
│  • Clear chat: C                                           │
│  • Save chat: S                                            │
│  • Load chat: L                                            │
│                                                             │
│  Chat History:                                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ Welcome to MarChat!                                    │ │
│  │ Type your message below...                             │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                             │
│  Message: [                                                ] │
│                                                             │
│  Press 'h' for help, 'q' to quit                          │
└─────────────────────────────────────────────────────────────┘`
}

func (p *Plugin) Name() string {
	return "marchat"
}

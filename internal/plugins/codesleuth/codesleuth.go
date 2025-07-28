package codesleuth

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
┌─ CodeSleuth ───────────────────────────────────────────────┐
│                                                             │
│  Static Analysis                                            │
│                                                             │
│  • Analyze current file: A                                 │
│  • Show IR diagram: I                                      │
│  • Find references: R                                      │
│  • Show call graph: C                                      │
│                                                             │
│  Current analysis:                                          │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ No analysis available                                  │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                             │
│  Press 'h' for help, 'q' to quit                          │
└─────────────────────────────────────────────────────────────┘`
}

func (p *Plugin) Name() string {
	return "codesleuth"
}

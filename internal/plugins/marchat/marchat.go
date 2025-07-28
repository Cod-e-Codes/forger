package marchat

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"forger/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

type Plugin struct {
	ctx           *types.Context
	serverRunning bool
	serverURL     string
	username      string
	theme         string
	messages      []Message
	input         string
	connected     bool
	errorMsg      string
}

type Message struct {
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}

func New(ctx *types.Context) types.Plugin {
	return &Plugin{
		ctx:       ctx,
		serverURL: "ws://localhost:9090/ws",
		username:  "ForgerUser",
		theme:     "patriot",
		messages:  []Message{},
	}
}

func (p *Plugin) Init() tea.Cmd {
	return p.checkServer
}

func (p *Plugin) checkServer() tea.Msg {
	// Check if marchat server is running
	cmd := exec.Command("marchat-client", "--help")
	if err := cmd.Run(); err != nil {
		return ServerCheckMsg{Available: false, Error: "marchat-client not found in PATH"}
	}
	return ServerCheckMsg{Available: true}
}

func (p *Plugin) Update(msg tea.Msg) (types.Plugin, tea.Cmd) {
	switch msg := msg.(type) {
	case ServerCheckMsg:
		p.serverRunning = msg.Available
		if !msg.Available {
			p.errorMsg = "MarChat client not found. Install from: https://github.com/Cod-e-Codes/marchat"
		}
		return p, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if p.input != "" {
				p.messages = append(p.messages, Message{
					Username:  p.username,
					Content:   p.input,
					Timestamp: time.Now(),
					Type:      "message",
				})
				p.input = ""
			}
		case "backspace":
			if len(p.input) > 0 {
				p.input = p.input[:len(p.input)-1]
			}
		case "ctrl+c":
			return p, tea.Quit
		default:
			if len(msg.String()) == 1 {
				p.input += msg.String()
			}
		}
	}
	return p, nil
}

func (p *Plugin) View() string {
	var sb strings.Builder

	sb.WriteString("┌─ MarChat ───────────────────────────────────────────────────┐\n")
	sb.WriteString("│                                                             │\n")

	if !p.serverRunning {
		sb.WriteString("│  ❌ MarChat Server Not Available                        │\n")
		sb.WriteString("│                                                             │\n")
		sb.WriteString("│  To use MarChat:                                       │\n")
		sb.WriteString("│  1. Install marchat:                                   │\n")
		sb.WriteString("│     go install github.com/Cod-e-Codes/marchat@latest   │\n")
		sb.WriteString("│  2. Start server: marchat-server                       │\n")
		sb.WriteString("│  3. Connect client: marchat-client                     │\n")
		sb.WriteString("│                                                             │\n")
		sb.WriteString("│  Server Status: " + p.getServerStatus() + "                    │\n")
	} else {
		sb.WriteString("│  ✅ MarChat Server Available                            │\n")
		sb.WriteString("│                                                             │\n")
		sb.WriteString("│  Chat History:                                          │\n")
		sb.WriteString("│  ┌─────────────────────────────────────────────────────┐ │\n")

		// Show last 5 messages
		start := len(p.messages) - 5
		if start < 0 {
			start = 0
		}
		for i := start; i < len(p.messages); i++ {
			msg := p.messages[i]
			timeStr := msg.Timestamp.Format("15:04")
			line := fmt.Sprintf("│  [%s] %s: %s", timeStr, msg.Username, msg.Content)
			if len(line) > 55 {
				line = line[:52] + "..."
			}
			sb.WriteString(fmt.Sprintf("│  %-55s │\n", line))
		}

		// Fill remaining space
		remaining := 5 - (len(p.messages) - start)
		for i := 0; i < remaining; i++ {
			sb.WriteString("│                                                         │\n")
		}

		sb.WriteString("│  └─────────────────────────────────────────────────────┘ │\n")
		sb.WriteString("│                                                             │\n")
		sb.WriteString(fmt.Sprintf("│  Message: [%-45s] │\n", p.input))
	}

	sb.WriteString("│                                                             │\n")
	sb.WriteString("│  Commands:                                                │\n")
	sb.WriteString("│  • Enter: Send message                                    │\n")
	sb.WriteString("│  • Ctrl+C: Quit                                           │\n")
	sb.WriteString("│  • Backspace: Edit message                                │\n")
	sb.WriteString("└─────────────────────────────────────────────────────────────┘")

	return sb.String()
}

func (p *Plugin) getServerStatus() string {
	if p.serverRunning {
		return "Running"
	}
	return "Not Found"
}

func (p *Plugin) Name() string {
	return "marchat"
}

type ServerCheckMsg struct {
	Available bool
	Error     string
}

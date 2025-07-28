package marchat

import (
	"fmt"
	"os"
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
	result        string // Add result field for command feedback
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
	return p.checkServer()
}

func (p *Plugin) checkServer() tea.Cmd {
	return func() tea.Msg {
		// Check if marchat server is running using full path
		marchatClientPath := "C:\\Users\\codyl\\go\\bin\\marchat-client.exe"

		// Debug: print the path being used
		fmt.Printf("DEBUG: Looking for marchat-client at: %s\n", marchatClientPath)

		// Debug: check if file exists
		if _, err := os.Stat(marchatClientPath); os.IsNotExist(err) {
			return ServerCheckMsg{Available: false, Error: fmt.Sprintf("marchat-client not found at: %s", marchatClientPath)}
		}

		cmd := exec.Command(marchatClientPath, "--help")
		if err := cmd.Run(); err != nil {
			return ServerCheckMsg{Available: false, Error: fmt.Sprintf("marchat-client failed to run: %v", err)}
		}
		return ServerCheckMsg{Available: true}
	}
}

func (p *Plugin) Update(msg tea.Msg) (types.Plugin, tea.Cmd) {
	switch msg := msg.(type) {
	case ServerCheckMsg:
		p.serverRunning = msg.Available
		p.errorMsg = msg.Error
		// Update connection status based on server availability
		p.connected = msg.Available
		return p, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if p.input != "" {
				// Send message logic would go here
				p.messages = append(p.messages, Message{Username: "You", Content: p.input, Timestamp: time.Now(), Type: "message"})
				if p.connected {
					p.result = "✅ Message sent: " + p.input
				} else {
					p.result = "❌ Not connected to server"
				}
				p.input = ""
			}
		case "backspace":
			if len(p.input) > 0 {
				p.input = p.input[:len(p.input)-1]
			}
		case "ctrl+c":
			return p, tea.Quit
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

		// Show connection status
		if p.connected {
			sb.WriteString("│  🔗 Connected to server                               │\n")
		} else {
			sb.WriteString("│  🔌 Disconnected from server                          │\n")
		}
		sb.WriteString("│                                                             │\n")

		// Show feedback if any
		if p.result != "" {
			sb.WriteString("│  Status: " + p.result + "                                    │\n")
			sb.WriteString("│                                                             │\n")
		}

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

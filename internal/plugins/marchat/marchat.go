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
	serverProcess *os.Process
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
	return p.startServer()
}

func (p *Plugin) startServer() tea.Cmd {
	return func() tea.Msg {
		// Check if we have the marchat server executable
		marchatServerPath := os.Getenv("GOPATH") + "\\bin\\marchat-server.exe"

		fmt.Printf("DEBUG: MarChat checking server at: %s\n", marchatServerPath)

		// Check if server executable exists
		if _, err := os.Stat(marchatServerPath); os.IsNotExist(err) {
			fmt.Printf("DEBUG: MarChat server not found: %v\n", err)
			return ServerCheckMsg{Available: false, Error: fmt.Sprintf("marchat-server not found at: %s", marchatServerPath)}
		}

		fmt.Printf("DEBUG: MarChat server found, starting...\n")

		// Start the server using the executable
		cmd := exec.Command(marchatServerPath, "-config", "server_config.json")
		// Don't redirect stdout/stderr so we can see any error messages

		if err := cmd.Start(); err != nil {
			fmt.Printf("DEBUG: MarChat failed to start server: %v\n", err)
			return ServerCheckMsg{Available: false, Error: fmt.Sprintf("failed to start marchat-server: %v", err)}
		}

		fmt.Printf("DEBUG: MarChat server started, PID: %d\n", cmd.Process.Pid)

		// Store the process so we can kill it later
		p.serverProcess = cmd.Process

		// Wait longer for server to start
		time.Sleep(5 * time.Second)

		// Check if server is responding by trying to connect
		marchatClientPath := os.Getenv("GOPATH") + "\\bin\\marchat-client.exe"
		testCmd := exec.Command(marchatClientPath, "-username", "ForgerUser", "-admin", "-admin-key", "forger-admin-key", "-server", "ws://localhost:9090/ws")

		fmt.Printf("DEBUG: MarChat testing server connection...\n")
		output, err := testCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("DEBUG: MarChat server test failed: %v\nOutput: %s\n", err, string(output))
			return ServerCheckMsg{Available: false, Error: fmt.Sprintf("server not responding: %v", err)}
		}

		fmt.Printf("DEBUG: MarChat server test output: %s\n", string(output))
		fmt.Printf("DEBUG: MarChat server is responding!\n")
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
				// Send message using marchat-client
				marchatClientPath := os.Getenv("GOPATH") + "\\bin\\marchat-client.exe"
				cmd := exec.Command(marchatClientPath, "-username", "ForgerUser", "-admin", "-admin-key", "forger-admin-key", "-server", "ws://localhost:9090/ws")
				cmd.Stdin = strings.NewReader(p.input + "\n")

				err := cmd.Run()
				if err != nil {
					p.result = "❌ Failed to send message: " + err.Error()
				} else {
					p.result = "✅ Message sent: " + p.input
					p.messages = append(p.messages, Message{Username: "You", Content: p.input, Timestamp: time.Now(), Type: "message"})
				}
				p.input = ""
			}
		case "backspace":
			if len(p.input) > 0 {
				p.input = p.input[:len(p.input)-1]
			}
		case "ctrl+c":
			// Clean up server process
			if p.serverProcess != nil {
				p.serverProcess.Kill()
			}
			return p, tea.Quit
		default:
			// Handle regular character input
			if len(msg.String()) == 1 && msg.String() != "tab" && msg.String() != "shift+tab" {
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

		// Show connection status
		if p.serverRunning {
			sb.WriteString("│  🔗 Server started and running                        │\n")
		} else {
			sb.WriteString("│  🔌 Server not running                               │\n")
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

package codesleuth

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"forger/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

type Plugin struct {
	ctx           *types.Context
	available     bool
	selectedIndex int
	result        string // Add result field for command feedback
}

func New(ctx *types.Context) types.Plugin {
	return &Plugin{
		ctx:           ctx,
		selectedIndex: 0,
	}
}

func (p *Plugin) Init() tea.Cmd {
	return p.checkAvailability()
}

func (p *Plugin) checkAvailability() tea.Cmd {
	return func() tea.Msg {
		codesleuthPath := "C:\\Users\\codyl\\go\\bin\\codesleuth.exe"

		if _, err := os.Stat(codesleuthPath); os.IsNotExist(err) {
			return AvailabilityMsg{Available: false, Error: fmt.Sprintf("codesleuth not found at: %s", codesleuthPath)}
		}

		cmd := exec.Command(codesleuthPath, "--help")
		if err := cmd.Run(); err != nil {
			return AvailabilityMsg{Available: false, Error: fmt.Sprintf("codesleuth failed to run: %v", err)}
		}
		return AvailabilityMsg{Available: true}
	}
}

func (p *Plugin) Update(msg tea.Msg) (types.Plugin, tea.Cmd) {
	switch msg := msg.(type) {
	case AvailabilityMsg:
		p.available = msg.Available
		return p, nil
	case CommandResultMsg:
		// Display command results
		if msg.Success {
			p.result = "✅ " + msg.Output
		} else {
			p.result = "❌ " + msg.Output
		}
		return p, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			return p, p.analyzeCurrentDirectory
		case "i":
			return p, p.showIRDiagram
		case "r":
			return p, p.findReferences
		case "g":
			return p, p.showCallGraph
		case "ctrl+c":
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p *Plugin) analyzeCurrentDirectory() tea.Msg {
	codesleuthPath := os.Getenv("GOPATH") + "\\bin\\codesleuth.exe"
	cmd := exec.Command(codesleuthPath, "analyze", ".")
	output, err := cmd.CombinedOutput() // Use CombinedOutput to get both stdout and stderr
	if err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("CodeSleuth only supports COBOL files. Current directory contains Go files.\nError: %v\nOutput: %s", err, string(output)),
		}
	}

	// Since CodeSleuth doesn't support JSON output, we'll just show the raw output
	return CommandResultMsg{
		Success: true,
		Output:  fmt.Sprintf("CodeSleuth analysis output:\n%s", string(output)),
	}
}

func (p *Plugin) showIRDiagram() tea.Msg {
	codesleuthPath := os.Getenv("GOPATH") + "\\bin\\codesleuth.exe"
	cmd := exec.Command(codesleuthPath, "analyze", ".", "--mermaid")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("Error generating IR diagram: %v\n%s", err, string(output)),
		}
	}

	return CommandResultMsg{
		Success: true,
		Output:  fmt.Sprintf("IR Diagram:\n%s", string(output)),
	}
}

func (p *Plugin) findReferences() tea.Msg {
	codesleuthPath := os.Getenv("GOPATH") + "\\bin\\codesleuth.exe"
	cmd := exec.Command(codesleuthPath, "analyze", ".", "--references")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("Error finding references: %v\n%s", err, string(output)),
		}
	}

	return CommandResultMsg{
		Success: true,
		Output:  fmt.Sprintf("References:\n%s", string(output)),
	}
}

func (p *Plugin) showCallGraph() tea.Msg {
	codesleuthPath := os.Getenv("GOPATH") + "\\bin\\codesleuth.exe"
	cmd := exec.Command(codesleuthPath, "analyze", ".", "--call-graph")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("Error generating call graph: %v\n%s", err, string(output)),
		}
	}

	return CommandResultMsg{
		Success: true,
		Output:  fmt.Sprintf("Call Graph:\n%s", string(output)),
	}
}

func (p *Plugin) View() string {
	var sb strings.Builder

	sb.WriteString("┌─ CodeSleuth ───────────────────────────────────────────────┐\n")
	sb.WriteString("│                                                             │\n")

	if !p.available {
		sb.WriteString("│  ❌ CodeSleuth Not Available                              │\n")
		sb.WriteString("│                                                             │\n")
		sb.WriteString("│  To use CodeSleuth:                                      │\n")
		sb.WriteString("│  1. Install codesleuth:                                  │\n")
		sb.WriteString("│     go install github.com/Cod-e-Codes/codesleuth@latest │\n")
		sb.WriteString("│  2. Build Rust components:                              │\n")
		sb.WriteString("│     cd parser && cargo build --release                  │\n")
		sb.WriteString("│     cd summarizer && cargo build --release              │\n")
		sb.WriteString("│  3. Analyze files: codesleuth analyze <path>           │\n")
		sb.WriteString("│                                                             │\n")
	} else {
		sb.WriteString("│  ✅ CodeSleuth Available                                  │\n")
		sb.WriteString("│                                                             │\n")

		// Show command results if any
		if p.result != "" {
			sb.WriteString("│  Result:                                                │\n")
			sb.WriteString("│  ┌─────────────────────────────────────────────────────┐ │\n")
			lines := strings.Split(p.result, "\n")
			for i, line := range lines {
				if i >= 3 { // Limit to 3 lines
					sb.WriteString("│  │ ... (truncated)                                    │ │\n")
					break
				}
				if len(line) > 55 {
					line = line[:52] + "..."
				}
				sb.WriteString(fmt.Sprintf("│  │ %-55s │ │\n", line))
			}
			sb.WriteString("│  └─────────────────────────────────────────────────────┘ │\n")
			sb.WriteString("│                                                             │\n")
		}

		sb.WriteString("│  Analysis Files:                                          │\n")
		sb.WriteString("│  ┌─────────────────────────────────────────────────────┐ │\n")

		sb.WriteString("│  │ No COBOL files analyzed                              │ │\n")
		sb.WriteString("│  │ Press 'A' to analyze current directory            │ │\n")
		sb.WriteString("│  │ (Note: CodeSleuth only supports COBOL files)      │ │\n")

		sb.WriteString("│  └─────────────────────────────────────────────────────┘ │\n")
		sb.WriteString("│                                                             │\n")
		sb.WriteString("│  Commands:                                                │\n")
		sb.WriteString("│  • A: Analyze current directory (COBOL files only)      │\n")
		sb.WriteString("│  • I: Show IR diagram                                    │\n")
		sb.WriteString("│  • R: Find references                                    │\n")
		sb.WriteString("│  • G: Show call graph                                    │\n")
	}

	sb.WriteString("│                                                             │\n")
	sb.WriteString("└─────────────────────────────────────────────────────────────┘")

	return sb.String()
}

func (p *Plugin) Name() string {
	return "codesleuth"
}

type AvailabilityMsg struct {
	Available bool
	Error     string
}

type CommandResultMsg struct {
	Success bool
	Output  string
}

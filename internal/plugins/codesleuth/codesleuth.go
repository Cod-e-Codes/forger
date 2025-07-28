package codesleuth

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"forger/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

type Plugin struct {
	ctx           *types.Context
	available     bool
	analysisFiles []AnalysisFile
	selectedIndex int
	output        string
	errorMsg      string
	status        string
}

type AnalysisFile struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Type       string `json:"type"`
	Lines      int    `json:"lines"`
	Functions  int    `json:"functions"`
	Variables  int    `json:"variables"`
	DeadCode   int    `json:"dead_code"`
	Complexity int    `json:"complexity"`
}

type AnalysisResult struct {
	ProgramName string         `json:"program_name"`
	Author      string         `json:"author"`
	Files       []AnalysisFile `json:"files"`
	Summary     string         `json:"summary"`
}

func New(ctx *types.Context) types.Plugin {
	return &Plugin{
		ctx:           ctx,
		selectedIndex: 0,
		analysisFiles: []AnalysisFile{},
	}
}

func (p *Plugin) Init() tea.Cmd {
	return p.checkAvailability
}

func (p *Plugin) checkAvailability() tea.Msg {
	cmd := exec.Command("codesleuth", "--help")
	if err := cmd.Run(); err != nil {
		return AvailabilityMsg{Available: false, Error: "codesleuth not found in PATH"}
	}
	return AvailabilityMsg{Available: true}
}

func (p *Plugin) Update(msg tea.Msg) (types.Plugin, tea.Cmd) {
	switch msg := msg.(type) {
	case AvailabilityMsg:
		p.available = msg.Available
		return p, nil
	case AnalysisResultMsg:
		p.analysisFiles = msg.Result.Files
		p.output = msg.Result.Summary
		return p, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if p.selectedIndex > 0 {
				p.selectedIndex--
			}
		case "down":
			if p.selectedIndex < len(p.analysisFiles)-1 {
				p.selectedIndex++
			}
		case "enter":
			if len(p.analysisFiles) > 0 && p.selectedIndex < len(p.analysisFiles) {
				return p, func() tea.Msg { return p.analyzeFile(p.analysisFiles[p.selectedIndex]) }
			}
		case "a":
			return p, p.analyzeCurrentDirectory
		case "i":
			if len(p.analysisFiles) > 0 {
				return p, func() tea.Msg { return p.showIRDiagram(p.analysisFiles[p.selectedIndex]) }
			}
		case "r":
			if len(p.analysisFiles) > 0 {
				return p, func() tea.Msg { return p.findReferences(p.analysisFiles[p.selectedIndex]) }
			}
		case "c":
			if len(p.analysisFiles) > 0 {
				return p, func() tea.Msg { return p.showCallGraph(p.analysisFiles[p.selectedIndex]) }
			}
		case "ctrl+c":
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p *Plugin) analyzeCurrentDirectory() tea.Msg {
	cmd := exec.Command("codesleuth", "analyze", ".", "--json")
	output, err := cmd.Output()
	if err != nil {
		return AnalysisResultMsg{
			Result: AnalysisResult{
				ProgramName: "Error",
				Author:      "Unknown",
				Files:       []AnalysisFile{},
				Summary:     fmt.Sprintf("Error analyzing directory: %v", err),
			},
		}
	}

	var result AnalysisResult
	if err := json.Unmarshal(output, &result); err != nil {
		return AnalysisResultMsg{
			Result: AnalysisResult{
				ProgramName: "Parse Error",
				Author:      "Unknown",
				Files:       []AnalysisFile{},
				Summary:     fmt.Sprintf("Error parsing analysis: %v", err),
			},
		}
	}

	return AnalysisResultMsg{Result: result}
}

func (p *Plugin) analyzeFile(file AnalysisFile) tea.Msg {
	cmd := exec.Command("codesleuth", "analyze", file.Path, "--json")
	output, err := cmd.Output()
	if err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("Error analyzing file %s: %v", file.Name, err),
		}
	}

	var result AnalysisResult
	if err := json.Unmarshal(output, &result); err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("Error parsing analysis for %s: %v", file.Name, err),
		}
	}

	return CommandResultMsg{
		Success: true,
		Output:  fmt.Sprintf("Analysis for %s:\n%s", file.Name, result.Summary),
	}
}

func (p *Plugin) showIRDiagram(file AnalysisFile) tea.Msg {
	cmd := exec.Command("codesleuth", "analyze", file.Path, "--mermaid")
	output, err := cmd.Output()
	if err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("Error generating IR diagram for %s: %v", file.Name, err),
		}
	}

	return CommandResultMsg{
		Success: true,
		Output:  fmt.Sprintf("IR Diagram for %s:\n%s", file.Name, output),
	}
}

func (p *Plugin) findReferences(file AnalysisFile) tea.Msg {
	cmd := exec.Command("codesleuth", "analyze", file.Path, "--references")
	output, err := cmd.Output()
	if err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("Error finding references for %s: %v", file.Name, err),
		}
	}

	return CommandResultMsg{
		Success: true,
		Output:  fmt.Sprintf("References for %s:\n%s", file.Name, output),
	}
}

func (p *Plugin) showCallGraph(file AnalysisFile) tea.Msg {
	cmd := exec.Command("codesleuth", "analyze", file.Path, "--call-graph")
	output, err := cmd.Output()
	if err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("Error generating call graph for %s: %v", file.Name, err),
		}
	}

	return CommandResultMsg{
		Success: true,
		Output:  fmt.Sprintf("Call Graph for %s:\n%s", file.Name, output),
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
		sb.WriteString("│  Analysis Files:                                          │\n")
		sb.WriteString("│  ┌─────────────────────────────────────────────────────┐ │\n")

		if len(p.analysisFiles) == 0 {
			sb.WriteString("│  │ No files analyzed                                  │ │\n")
			sb.WriteString("│  │ Press 'A' to analyze current directory            │ │\n")
		} else {
			for i, file := range p.analysisFiles {
				prefix := "  "
				if i == p.selectedIndex {
					prefix = "> "
				}
				line := fmt.Sprintf("│  %s%s (%s) - %d lines", prefix, file.Name, file.Type, file.Lines)
				if len(line) > 55 {
					line = line[:52] + "..."
				}
				sb.WriteString(fmt.Sprintf("│  %-55s │\n", line))
			}
		}

		sb.WriteString("│  └─────────────────────────────────────────────────────┘ │\n")
		sb.WriteString("│                                                             │\n")
		sb.WriteString("│  Commands:                                                │\n")
		sb.WriteString("│  • A: Analyze current directory                          │\n")
		sb.WriteString("│  • I: Show IR diagram                                    │\n")
		sb.WriteString("│  • R: Find references                                    │\n")
		sb.WriteString("│  • C: Show call graph                                    │\n")
		sb.WriteString("│  • ↑/↓: Navigate files                                   │\n")
		sb.WriteString("│  • Enter: Analyze selected file                          │\n")
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

type AnalysisResultMsg struct {
	Result AnalysisResult
}

type CommandResultMsg struct {
	Success bool
	Output  string
}

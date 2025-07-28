package ignoregrets

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
	available     bool
	snapshots     []Snapshot
	selectedIndex int
	output        string
	errorMsg      string
	status        string
}

type Snapshot struct {
	Commit    string    `json:"commit"`
	Timestamp time.Time `json:"timestamp"`
	Index     int       `json:"index"`
	FileCount int       `json:"file_count"`
}

func New(ctx *types.Context) types.Plugin {
	return &Plugin{
		ctx:           ctx,
		selectedIndex: 0,
		snapshots:     []Snapshot{},
	}
}

func (p *Plugin) Init() tea.Cmd {
	return p.checkAvailability
}

func (p *Plugin) checkAvailability() tea.Msg {
	ignoregretsPath := os.Getenv("GOPATH") + "\\bin\\ignoregrets.exe"
	cmd := exec.Command(ignoregretsPath, "--help")
	if err := cmd.Run(); err != nil {
		return AvailabilityMsg{Available: false, Error: "ignoregrets not found"}
	}
	return AvailabilityMsg{Available: true}
}

func (p *Plugin) Update(msg tea.Msg) (types.Plugin, tea.Cmd) {
	switch msg := msg.(type) {
	case AvailabilityMsg:
		p.available = msg.Available
		if msg.Available {
			return p, p.listSnapshots
		}
		return p, nil
	case SnapshotsMsg:
		p.snapshots = msg.Snapshots
		return p, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if p.selectedIndex > 0 {
				p.selectedIndex--
			}
		case "down":
			if p.selectedIndex < len(p.snapshots)-1 {
				p.selectedIndex++
			}
		case "enter":
			if len(p.snapshots) > 0 && p.selectedIndex < len(p.snapshots) {
				return p, func() tea.Msg { return p.restoreSnapshot(p.snapshots[p.selectedIndex]) }
			}
		case "s":
			return p, p.createSnapshot
		case "r":
			if len(p.snapshots) > 0 {
				return p, func() tea.Msg { return p.restoreSnapshot(p.snapshots[p.selectedIndex]) }
			}
		case "l":
			return p, p.listSnapshots
		case "d":
			if len(p.snapshots) > 0 {
				return p, func() tea.Msg { return p.deleteSnapshot(p.snapshots[p.selectedIndex]) }
			}
		case "ctrl+c":
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p *Plugin) createSnapshot() tea.Msg {
	ignoregretsPath := os.Getenv("GOPATH") + "\\bin\\ignoregrets.exe"
	cmd := exec.Command(ignoregretsPath, "snapshot")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("Error creating snapshot: %v\n%s", err, output),
		}
	}
	return CommandResultMsg{
		Success: true,
		Output:  fmt.Sprintf("Snapshot created successfully:\n%s", output),
	}
}

func (p *Plugin) listSnapshots() tea.Msg {
	ignoregretsPath := os.Getenv("GOPATH") + "\\bin\\ignoregrets.exe"
	cmd := exec.Command(ignoregretsPath, "list")
	output, err := cmd.Output()
	if err != nil {
		return SnapshotsMsg{Snapshots: []Snapshot{}}
	}

	// Parse the output to extract snapshots
	snapshots := p.parseSnapshots(string(output))
	return SnapshotsMsg{Snapshots: snapshots}
}

func (p *Plugin) restoreSnapshot(snapshot Snapshot) tea.Msg {
	ignoregretsPath := os.Getenv("GOPATH") + "\\bin\\ignoregrets.exe"
	cmd := exec.Command(ignoregretsPath, "restore", "--dry-run")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("Error previewing restore: %v\n%s", err, output),
		}
	}
	return CommandResultMsg{
		Success: true,
		Output:  fmt.Sprintf("Restore preview for %s:\n%s", snapshot.Commit[:8], output),
	}
}

func (p *Plugin) deleteSnapshot(snapshot Snapshot) tea.Msg {
	ignoregretsPath := os.Getenv("GOPATH") + "\\bin\\ignoregrets.exe"
	cmd := exec.Command(ignoregretsPath, "prune", "--retention", "0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return CommandResultMsg{
			Success: false,
			Output:  fmt.Sprintf("Error pruning snapshots: %v\n%s", err, output),
		}
	}
	return CommandResultMsg{
		Success: true,
		Output:  fmt.Sprintf("Snapshots pruned (including %s):\n%s", snapshot.Commit[:8], output),
	}
}

func (p *Plugin) parseSnapshots(output string) []Snapshot {
	var snapshots []Snapshot
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.Contains(line, "Commit:") {
			// Parse commit line
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				snapshot := Snapshot{
					Commit:    parts[1],
					Timestamp: time.Now(),
					Index:     len(snapshots),
					FileCount: 0,
				}
				snapshots = append(snapshots, snapshot)
			}
		}
	}

	return snapshots
}

func (p *Plugin) View() string {
	var sb strings.Builder

	sb.WriteString("┌─ IgnoreGrets ──────────────────────────────────────────────┐\n")
	sb.WriteString("│                                                             │\n")

	if !p.available {
		sb.WriteString("│  ❌ IgnoreGrets Not Available                            │\n")
		sb.WriteString("│                                                             │\n")
		sb.WriteString("│  To use IgnoreGrets:                                    │\n")
		sb.WriteString("│  1. Install ignoregrets:                                │\n")
		sb.WriteString("│     go install github.com/Cod-e-Codes/ignoregrets@latest│\n")
		sb.WriteString("│  2. Initialize in repo: ignoregrets init                │\n")
		sb.WriteString("│  3. Create snapshots: ignoregrets snapshot              │\n")
		sb.WriteString("│                                                             │\n")
	} else {
		sb.WriteString("│  ✅ IgnoreGrets Available                                │\n")
		sb.WriteString("│                                                             │\n")
		sb.WriteString("│  Snapshots:                                              │\n")
		sb.WriteString("│  ┌─────────────────────────────────────────────────────┐ │\n")

		if len(p.snapshots) == 0 {
			sb.WriteString("│  │ No snapshots available                              │ │\n")
			sb.WriteString("│  │ Run 'ignoregrets snapshot' to create one           │ │\n")
		} else {
			for i, snapshot := range p.snapshots {
				prefix := "  "
				if i == p.selectedIndex {
					prefix = "> "
				}
				timeStr := snapshot.Timestamp.Format("2006-01-02 15:04")
				line := fmt.Sprintf("│  %s%s (%d files) - %s", prefix, snapshot.Commit[:8], snapshot.FileCount, timeStr)
				if len(line) > 55 {
					line = line[:52] + "..."
				}
				sb.WriteString(fmt.Sprintf("│  %-55s │\n", line))
			}
		}

		sb.WriteString("│  └─────────────────────────────────────────────────────┘ │\n")
		sb.WriteString("│                                                             │\n")
		sb.WriteString("│  Commands:                                                │\n")
		sb.WriteString("│  • S: Create snapshot                                    │\n")
		sb.WriteString("│  • R: Restore snapshot (preview)                         │\n")
		sb.WriteString("│  • D: Delete old snapshots                               │\n")
		sb.WriteString("│  • L: Refresh list                                       │\n")
		sb.WriteString("│  • ↑/↓: Navigate snapshots                               │\n")
		sb.WriteString("│  • Enter: Restore selected snapshot                      │\n")
	}

	sb.WriteString("│                                                             │\n")
	sb.WriteString("└─────────────────────────────────────────────────────────────┘")

	return sb.String()
}

func (p *Plugin) Name() string {
	return "ignoregrets"
}

type AvailabilityMsg struct {
	Available bool
	Error     string
}

type SnapshotsMsg struct {
	Snapshots []Snapshot
}

type CommandResultMsg struct {
	Success bool
	Output  string
}

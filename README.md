# Forger

[![Go](https://img.shields.io/badge/go-1.21%2B-blue)](https://go.dev) [![MIT License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Forger is a terminal-native developer dashboard built with Go, [Bubble Tea](https://github.com/charmbracelet/bubbletea), and [Lip Gloss](https://github.com/charmbracelet/lipgloss). It integrates your favorite CLI tools into a unified TUI for static analysis, code navigation, environment inspection, and more.

## What is Forger?

Forger is a **terminal-native developer toolkit** that provides a unified interface for multiple development tools. Think of it as a dashboard for your CLI tools - instead of switching between different terminal windows or running commands separately, Forger brings everything together in one interactive terminal interface.

## Features

- **Plugin Architecture**: Each tool is implemented as a plugin with a consistent interface
- **Terminal-First UX**: Built entirely for the terminal using Bubble Tea and Lip Gloss
- **Keyboard-Driven**: Navigate between plugins using arrow keys and keyboard shortcuts
- **Extensible**: Easy to add new plugins or modify existing ones
- **Fast & Lightweight**: No GUI frameworks or heavy dependencies
- **Real Tool Integration**: Connects to actual CLI tools (marchat, ignoregrets, codesleuth)

## Current Plugins

### IgnoreGrets ✅ **Integrated**
- **Purpose**: Snapshot management and Git workflow tools
- **Features**: Create, list, restore, and delete snapshots
- **Integration**: Direct CLI wrapper with real ignoregrets commands
- **Use Case**: Managing code states and quick rollbacks
- **Installation**: `go install github.com/Cod-e-Codes/ignoregrets@latest`

### CodeSleuth ✅ **Integrated**  
- **Purpose**: Static code analysis and IR visualization
- **Features**: Analyze files, show IR diagrams, find references, call graphs
- **Integration**: CLI wrapper with JSON output parsing
- **Use Case**: Understanding code structure and dependencies
- **Installation**: `go install github.com/Cod-e-Codes/codesleuth@latest`

### MarChat ✅ **Integrated**
- **Purpose**: Terminal-based chat interface
- **Features**: Send messages, save/load chat history, clear conversations
- **Integration**: Direct integration with marchat server/client
- **Use Case**: Developer communication and note-taking
- **Installation**: `go install github.com/Cod-e-Codes/marchat@latest`

## Quick Start

1. **Build the project**:
   ```bash
   go build ./cmd/forger
   ```

2. **Install required tools** (optional, for full functionality):
   ```bash
   # Install marchat for chat functionality
   go install github.com/Cod-e-Codes/marchat@latest
   
   # Install ignoregrets for snapshot management
   go install github.com/Cod-e-Codes/ignoregrets@latest
   
   # Install codesleuth for code analysis
   go install github.com/Cod-e-Codes/codesleuth@latest
   ```

3. **Run Forger**:
   ```bash
   ./forger
   ```

4. **Navigate the interface**:
   - Use **Up/Down arrows** to switch between plugins
   - Press **'c'** to open MarChat overlay
   - Press **'q'** or **Ctrl+C** to quit
   - Press **'esc'** to close overlays

## Plugin-Specific Controls

### IgnoreGrets
- **S**: Create snapshot
- **R**: Restore snapshot (preview)
- **D**: Delete old snapshots
- **L**: Refresh list
- **↑/↓**: Navigate snapshots
- **Enter**: Restore selected snapshot

### CodeSleuth
- **A**: Analyze current directory
- **I**: Show IR diagram
- **R**: Find references
- **C**: Show call graph
- **↑/↓**: Navigate files
- **Enter**: Analyze selected file

### MarChat
- **Enter**: Send message
- **Backspace**: Edit message
- **Ctrl+C**: Quit

## Configuration

Forger uses `forger.json` for configuration:

```json
{
  "default": "ignoregrets",
  "enabled": [
    "ignoregrets",
    "codesleuth", 
    "marchat"
  ]
}
```

- `default`: The plugin to show when Forger starts
- `enabled`: List of plugins to load

## Architecture

```
forger/
├── cmd/forger/          # Main application entry point
├── internal/
│   ├── core/           # Core runtime and plugin management
│   ├── types/          # Shared interfaces and types
│   └── plugins/        # Individual plugin implementations
│       ├── ignoregrets/ # Git snapshot management
│       ├── codesleuth/  # Code analysis
│       └── marchat/     # Terminal chat
└── forger.json         # Configuration file
```

## Development

### Adding a New Plugin

1. Create a new directory in `internal/plugins/`
2. Implement the `Plugin` interface:
   ```go
   type Plugin interface {
       Init() tea.Cmd
       Update(msg tea.Msg) (Plugin, tea.Cmd)
       View() string
       Name() string
   }
   ```
3. Add the plugin to the registry in `internal/core/registry.go`
4. Update `forger.json` to enable the plugin

### Building

```bash
go build ./cmd/forger
```

## Status

This is a **working release** with real integrations:

- ✅ **IgnoreGrets**: Full CLI integration with snapshot management
- ✅ **CodeSleuth**: Code analysis with JSON output parsing
- ✅ **MarChat**: Chat interface with server detection
- 🔄 **Future**: Additional plugins (ascii-colorizer, parsec, etc.)

## Tool Dependencies

For full functionality, install these tools:

```bash
# Required for chat functionality
go install github.com/Cod-e-Codes/marchat@latest

# Required for snapshot management
go install github.com/Cod-e-Codes/ignoregrets@latest

# Required for code analysis
go install github.com/Cod-e-Codes/codesleuth@latest
```

## Future Enhancements

- Integration with additional tools (ascii-colorizer, parsec, etc.)
- Plugin registry and installation system
- Git-aware workspace detection
- Custom dashboards and layouts
- Real-time updates and notifications
- Plugin configuration management

## License

MIT License - see [LICENSE](LICENSE) for details.

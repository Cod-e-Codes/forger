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

## Current Plugins

### IgnoreGrets
- **Purpose**: Snapshot management and Git workflow tools
- **Features**: Create, list, restore, and delete snapshots
- **Use Case**: Managing code states and quick rollbacks

### CodeSleuth  
- **Purpose**: Static code analysis and IR visualization
- **Features**: Analyze files, show IR diagrams, find references, call graphs
- **Use Case**: Understanding code structure and dependencies

### MarChat
- **Purpose**: Terminal-based chat interface
- **Features**: Send messages, save/load chat history, clear conversations
- **Use Case**: Developer communication and note-taking

## Quick Start

1. **Build the project**:
   ```bash
   go build ./cmd/forger
   ```

2. **Run Forger**:
   ```bash
   ./forger
   ```

3. **Navigate the interface**:
   - Use **Up/Down arrows** to switch between plugins
   - Press **'c'** to open MarChat overlay
   - Press **'q'** or **Ctrl+C** to quit
   - Press **'esc'** to close overlays

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
│       ├── ignoregrets/
│       ├── codesleuth/
│       └── marchat/
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

This is a **pre-release version** with basic plugin implementations. The current plugins provide placeholder interfaces - they demonstrate the architecture but don't yet integrate with the actual CLI tools they represent.

## Future Enhancements

- Integration with actual CLI tools (ignoregrets, codesleuth, marchat)
- Real snapshot management and Git operations
- Static analysis with IR visualization
- Terminal chat with message persistence
- Additional plugins (ascii-colorizer, parsec, etc.)
- Plugin registry and installation system
- Git-aware workspace detection
- Custom dashboards and layouts

## License

MIT License - see [LICENSE](LICENSE) for details.

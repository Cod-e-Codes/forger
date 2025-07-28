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

### IgnoreGrets ✅ **Fully Integrated**
- **Purpose**: Snapshot management and Git workflow tools
- **Features**: Create, list, restore, and delete snapshots
- **Integration**: Direct CLI wrapper with real ignoregrets commands
- **Use Case**: Managing code states and quick rollbacks
- **Status**: ✅ **Working** - Available and functional

### CodeSleuth ✅ **Fully Integrated**  
- **Purpose**: Static code analysis and IR visualization
- **Features**: Analyze files, show IR diagrams, find references, call graphs
- **Integration**: CLI wrapper with JSON output parsing
- **Use Case**: Understanding code structure and dependencies
- **Status**: ✅ **Working** - Available and functional (COBOL files only)

### MarChat ⚠️ **Partially Integrated**
- **Purpose**: Terminal-based chat interface
- **Features**: Send messages, save/load chat history, clear conversations
- **Integration**: Direct integration with marchat server/client
- **Use Case**: Developer communication and note-taking
- **Status**: ⚠️ **Server Configuration Required** - Server starts but client connection needs manual setup
- **Note**: Requires `server_config.json` file with admin credentials

## Quick Start

### 1. Build Forger
```bash
go build ./cmd/forger
```

### 2. Install Required CLI Tools

**Important**: These tools must be installed and available in your `GOPATH/bin` directory for Forger to detect them properly.

#### Install MarChat
```bash
# Clone the repository
git clone https://github.com/Cod-e-Codes/marchat.git temp-marchat
cd temp-marchat

# Build the client and server
go build ./client
go build ./server

# Copy executables to GOPATH/bin
copy client.exe $env:GOPATH\bin\marchat-client.exe
copy server.exe $env:GOPATH\bin\marchat-server.exe

# Clean up
cd ..
Remove-Item -Recurse -Force temp-marchat
```

#### Install IgnoreGrets
```bash
# Clone the repository
git clone https://github.com/Cod-e-Codes/ignoregrets.git temp-ignoregrets
cd temp-ignoregrets

# Build the executable
go build .

# Copy executable to GOPATH/bin
copy ignoregrets.exe $env:GOPATH\bin\ignoregrets.exe

# Clean up
cd ..
Remove-Item -Recurse -Force temp-ignoregrets
```

#### Install CodeSleuth
```bash
# Clone the repository
git clone https://github.com/Cod-e-Codes/codesleuth.git temp-codesleuth
cd temp-codesleuth

# Build the executable
go build ./cmd

# Copy executable to GOPATH/bin
copy cmd.exe $env:GOPATH\bin\codesleuth.exe

# Clean up
cd ..
Remove-Item -Recurse -Force temp-codesleuth
```

### 3. Configure MarChat (Required)

Create a `server_config.json` file in the Forger directory:

```json
{
  "port": 9090,
  "admin_key": "forger-admin-key",
  "theme": "patriot",
  "admins": ["ForgerUser"]
}
```

### 4. Verify Installation
```bash
# Test that all executables are available
& "$env:GOPATH\bin\marchat-client.exe" --help
& "$env:GOPATH\bin\ignoregrets.exe" --help
& "$env:GOPATH\bin\codesleuth.exe" --help
```

### 5. Run Forger
```bash
./forger
```

### 6. Navigate the Interface
- Use **Tab** to switch between plugins
- Use **Shift+Tab** to switch backwards between plugins
- Press **'c'** to open MarChat overlay
- Press **'q'** or **Ctrl+C** to quit
- Press **'esc'** to close overlays

## Plugin-Specific Controls

### IgnoreGrets
- **S**: Create snapshot
- **R**: Restore snapshot (preview)
- **D**: Delete old snapshots
- **L**: Refresh list
- **↑/↓**: Navigate snapshots (when plugin is active)
- **Enter**: Restore selected snapshot

### CodeSleuth
- **A**: Analyze current directory (COBOL files only)
- **I**: Show IR diagram
- **R**: Find references
- **G**: Show call graph
- **↑/↓**: Navigate files (when plugin is active)
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

## Troubleshooting

### Plugins Not Available
If plugins show as "Not Available":

1. **Check GOPATH**: Ensure your `GOPATH` is set correctly
   ```bash
   echo $env:GOPATH
   # Should show: C:\Users\codyl\go
   ```

2. **Verify Executables**: Check that executables exist in `GOPATH/bin`
   ```bash
   Test-Path "$env:GOPATH\bin\marchat-client.exe"
   Test-Path "$env:GOPATH\bin\ignoregrets.exe"
   Test-Path "$env:GOPATH\bin\codesleuth.exe"
   ```

3. **Test Executables**: Try running them directly
   ```bash
   & "$env:GOPATH\bin\marchat-client.exe" --help
   & "$env:GOPATH\bin\ignoregrets.exe" --help
   & "$env:GOPATH\bin\codesleuth.exe" --help
   ```

4. **Rebuild Forger**: After installing tools, rebuild Forger
   ```bash
   go build ./cmd/forger
   ```

### MarChat Issues
- **Server won't start**: Ensure `server_config.json` exists with proper admin configuration
- **Client can't connect**: Verify server is running on port 9090
- **Admin authentication**: Use `ForgerUser` as username with admin key `forger-admin-key`

### Common Issues

- **"Plugin not found"**: Ensure the plugin is listed in `forger.json` under `enabled`
- **"Executable not found"**: Verify the tool was built and copied to `GOPATH/bin` correctly
- **"Permission denied"**: Run PowerShell as Administrator if needed
- **CodeSleuth errors**: Remember that CodeSleuth only supports COBOL files currently

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
├── forger.json         # Configuration file
└── server_config.json  # MarChat server configuration
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
- ✅ **CodeSleuth**: Code analysis with JSON output parsing (COBOL files only)
- ⚠️ **MarChat**: Chat interface with server auto-start (requires manual client setup)
- 🔄 **Future**: Additional plugins (ascii-colorizer, parsec, etc.)

## Tool Dependencies

For full functionality, these tools must be installed in `GOPATH/bin`:

- **MarChat**: `marchat-client.exe` and `marchat-server.exe` (requires `server_config.json`)
- **IgnoreGrets**: `ignoregrets.exe`
- **CodeSleuth**: `codesleuth.exe`

## Future Enhancements

- Integration with additional tools (ascii-colorizer, parsec, etc.)
- Plugin registry and installation system
- Git-aware workspace detection
- Custom dashboards and layouts
- Real-time updates and notifications
- Plugin configuration management
- Improved MarChat integration with automatic client connection

## License

MIT License - see [LICENSE](LICENSE) for details.

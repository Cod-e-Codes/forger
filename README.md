# Forger

![Forger TUI Screenshot](docs/screenshot.png)

Forger is a terminal-native developer dashboard and plugin toolkit built in Go, using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss). It integrates multiple CLI tools into a unified, high‑performance TUI environment for static analysis, code navigation, environment inspection, and more.

## Quick Start

```bash
git clone https://github.com/yourusername/forger.git
cd forger
go build ./cmd/forger
./forger
```

## Features

* Plugin-based architecture with live plugin switching
* First-class support for tools like `marchat`, `ignoregrets`, and `codesleuth`
* Sidebar navigation with descriptions and previews
* Overlay toggles for help, logs, and plugin actions
* Configurable via `forger.json`
* Extensible with minimal plugin boilerplate

## Project Structure

```
forger/
├── cmd/forger/           # Main entrypoint
│   └── main.go
├── internal/core/        # Core application logic
│   ├── context.go        # Shared application state
│   ├── log.go            # Logging and error capture
│   ├── messages.go       # Message definitions
│   ├── model.go          # Main Bubble Tea model
│   └── plugin.go         # Plugin interface and manager
├── internal/plugins/     # Built-in plugins
│   ├── ignoregrets/
│   ├── codesleuth/
│   └── marchat/
├── docs/                 # Documentation and assets (e.g., screenshots)
│   └── screenshot.png
└── forger.json           # Configuration file
```

## Installation

Requires Go 1.20+

```bash
git clone https://github.com/yourusername/forger.git
cd forger
go build ./cmd/forger
```

Optional dependencies (for latest features):

```bash
go install github.com/charmbracelet/bubbletea@latest
go install github.com/charmbracelet/lipgloss@latest
```

## Configuration

Forger loads plugins from a `forger.json` file in the working directory. The JSON schema includes an optional default plugin and an enabled list.

```json
{
  "default": "ignoregrets",
  "enabled": [
    "marchat",
    "ignoregrets",
    "codesleuth"
  ]
}
```

* `default` (optional) — plugin selected at startup
* `enabled`         — list of plugin names to load

## Usage

```bash
./forger
```

### Keybindings

| Key       | Action                          |
| --------- | ------------------------------- |
| Up / Down | Select previous/next plugin     |
| Tab       | Next plugin                     |
| c         | Toggle chat overlay (`marchat`) |
| Esc       | Close overlay                   |
| q, Ctrl+C | Quit Forger                     |

## Plugin Development

Plugins are Go packages placed under `internal/plugins/<name>` and must implement the following interface:

```go
type Plugin interface {
  Name() string           // Unique identifier
  Description() string    // Short summary for sidebar or help overlay
  Init() tea.Cmd          // Initialization command
  Update(msg tea.Msg) (Plugin, tea.Cmd)
  View() string           // Render output
}
```

To create a new plugin:

1. Create a folder under `internal/plugins/`, e.g., `internal/plugins/myplugin/`.
2. Implement the `Plugin` interface in `plugin.go`.
3. Register the plugin in `internal/core/registry.go` by adding a factory entry.
4. Add the plugin name to `enabled` in `forger.json`.

See `internal/plugins/marchat` for a complete example.

## Troubleshooting

* **No plugins loaded**: Verify that `forger.json` exists and lists valid names for `enabled`.
* **Plugin not found**: Ensure the plugin directory name matches the entry in `forger.json` and is registered in `registry.go`.
* **Visual glitches**: Resize the terminal or use a compatible emulator.
* **Verbose logging**: Run with `FORGER_LOG=debug forger` to enable detailed output.

## License

This project is licensed under the [MIT License](./LICENSE).

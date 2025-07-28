# Forger

[![Go](https://img.shields.io/badge/go-1.20%2B-blue)](https://go.dev) [![MIT License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Forger is a terminal-native developer dashboard built with Go, [Bubble Tea](https://github.com/charmbracelet/bubbletea), and [Lip Gloss](https://github.com/charmbracelet/lipgloss). It integrates your favorite CLI tools into a unified TUI for static analysis, code navigation, environment inspection, and more.

---

## Features

* **Plugin-based architecture**: Add, remove, or switch between tools with a few keystrokes.
* **Live plugin switching**: Use arrow keys to toggle between plugins at runtime.
* **Built-in support** for:

  * [`marchat`](https://github.com/Cod-e-Codes/marchat): lightweight terminal chat
  * [`ignoregrets`](https://github.com/Cod-e-Codes/ignoregrets): snapshot/restore Git-ignored files
  * [`codesleuth`](https://github.com/Cod-e-Codes/codesleuth): COBOL static analyzer
* **Lightweight config system**: Define defaults and enable/disable plugins via `forger.json`
* **Minimal plugin boilerplate**: Implement a simple `Plugin` interface and register it

> Overlay support is currently limited to `marchat`. Sidebar descriptions and help/log overlays are planned but not yet implemented.

---

## Quick Start

```bash
git clone https://github.com/Cod-e-Codes/forger.git
cd forger
go build -o forger ./cmd/forger
./forger
```

---

## Project Structure

```text
forger/
├── cmd/forger/             # CLI entrypoint
│   └── main.go
├── internal/
│   ├── core/               # Core TUI logic
│   │   ├── model.go        # Main Bubble Tea model
│   │   ├── plugin.go       # Plugin interface
│   │   ├── registry.go     # Plugin registration and lookup
│   │   ├── context.go      # Shared context passed to plugins
│   │   ├── messages.go     # Internal message types
│   │   └── log.go          # Logging helpers
│   └── plugins/            # Plugin implementations
│       ├── marchat/        # marchat plugin
│       │   └── marchat.go
│       ├── ignoregrets/    # ignoregrets plugin
│       │   └── ignoregrets.go
│       └── codesleuth/     # codesleuth plugin
│           └── codesleuth.go
├── go.mod
├── go.sum
├── forger.json             # User config file
├── LICENSE
└── docs/
    └── screenshot.png      # TUI preview
```

---

## Installation

Requires Go 1.20 or later.

```bash
go install github.com/Cod-e-Codes/forger/cmd/forger@latest
```

Optional dependencies (already in `go.mod`):

```bash
go get github.com/charmbracelet/bubbletea@latest
```

---

## Configuration

Create a `forger.json` file in the root directory or pass it as a flag.

```json
{
  "default": "marchat",
  "enabled": ["marchat", "ignoregrets", "codesleuth"]
}
```

* `default`: the plugin to launch initially
* `enabled`: list of plugin names to show in the sidebar

---

## Usage

```bash
./forger -config forger.json
```

### Keybindings

| Key            | Action                |
| -------------- | --------------------- |
| `↑` / `↓`      | Navigate plugins      |
| `c`            | Toggle plugin overlay |
| `esc`          | Close overlay         |
| `q` / `Ctrl+C` | Quit Forger           |

> `Tab` support for navigation is planned but not yet implemented.

---

## Plugin Development

To add your own plugin:

1. Implement the `Plugin` interface:

```go
// internal/core/plugin.go
type Plugin interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Plugin, tea.Cmd)
    View() string
    Name() string
}
```

2. Add a factory function to instantiate your plugin:

```go
package yourplugin

import (
    "forger/internal/core"
    tea "github.com/charmbracelet/bubbletea"
)

type Plugin struct{}

func New(ctx *core.Context) core.Plugin {
    return &Plugin{}
}

func (p *Plugin) Init() tea.Cmd                            { return nil }
func (p *Plugin) Update(msg tea.Msg) (core.Plugin, tea.Cmd) { return p, nil }
func (p *Plugin) View() string                             { return "Hello from your plugin!" }
func (p *Plugin) Name() string                             { return "yourplugin" }
```

3. Register it in `internal/core/registry.go`:

```go
var availablePlugins = map[string]PluginFactory{
    "yourplugin": func(ctx *Context) Plugin {
        return yourplugin.New(ctx)
    },
}
```

4. Add it to `forger.json` under `enabled`.

---

## Troubleshooting

* Plugin not showing? Check your `forger.json`.
* Overlay not toggling? Only `marchat` currently implements overlays.
* Errors at startup? Check stderr and the UI for plugin load messages.

---

## License

[MIT](LICENSE)

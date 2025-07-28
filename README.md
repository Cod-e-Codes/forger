# Forger

Forger is a terminal-native, plugin‑driven developer dashboard built with Go, Bubble Tea and Lip Gloss. It provides a unified interface for your CLI tools, allowing you to switch between tasks—code analysis, snapshots, chat and more—without leaving the terminal.

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Configuration](#configuration)
4. [Usage](#usage)
5. [Plugins](#plugins)
6. [Keybindings](#keybindings)
7. [Plugin Development](#plugin-development)
8. [Contributing](#contributing)
9. [License](#license)

---

## Features

* **Modular Plugin System**
  Dynamically load and configure plugins for snapshots, code analysis, chat overlays, and more.

* **Deterministic Navigation**
  Sidebar lists plugins in alphabetical order; cycle with up/down arrows or tab.

* **Overlay Support**
  Toggle a chat overlay (or other modal) without losing context.

* **Shared Context**
  Plugins share mutable state for seamless data exchange.

* **Error Reporting**
  Load errors and runtime issues are displayed in‑app or logged to stderr.

---

## Installation

1. Clone the repository

   ```bash
   git clone https://github.com/username/forger.git
   cd forger
   ```

2. Build the binary

   ```bash
   go build -o forger ./cmd/forger
   ```

3. Ensure your tools are installed and in `PATH` (e.g., Go, Bubble Tea dependencies).

---

## Configuration

Forger uses a simple JSON file (`forger.json`) in the working directory:

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

* `default`  – the plugin to select on startup
* `enabled`  – list of plugin names to load

---

## Usage

Run Forger from your terminal:

```bash
./forger
```

Navigate the interface:

* **Up/Down**: cycle through plugins
* **c**       : toggle chat overlay (`marchat` plugin)
* **Esc**     : close overlay
* **q / Ctrl+C**: quit Forger

Each plugin defines its own commands and views once activated.

---

## Plugins

Plugins live under `internal/plugins` (or your preferred layout) and must implement the `core.Plugin` interface:

```go
type Plugin interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Plugin, tea.Cmd)
    View() string
    Name() string
}
```

Included plugins:

* **ignoregrets** — view and restore Git‑ignored file snapshots
* **codesleuth**  — visualize COBOL IR and call/control flow graphs
* **marchat**     — chat overlay for your terminal chat application

---

## Keybindings

Forger reserves the following keys by default:

| Key       | Action                      |
| --------- | --------------------------- |
| ↑ / ↓     | Select previous/next plugin |
| c         | Toggle `marchat` overlay    |
| Esc       | Close overlay               |
| q, Ctrl+C | Quit Forger                 |

Individual plugins may introduce additional keybindings once active.

---

## Plugin Development

1. Create a new directory under `internal/plugins/`
2. Implement the `Plugin` interface
3. Register the plugin in `internal/core/registry.go`
4. Update `forger.json` and add your plugin name to `enabled`

Example skeleton:

```go
package myplugin

import tea "github.com/charmbracelet/bubbletea"

func New(ctx *core.Context) core.Plugin {
    return &Plugin{}
}

type Plugin struct {}

func (p *Plugin) Init() tea.Cmd { return nil }

func (p *Plugin) Update(msg tea.Msg) (core.Plugin, tea.Cmd) {
    return p, nil
}

func (p *Plugin) View() string {
    return "MyPlugin View"
}

func (p *Plugin) Name() string {
    return "myplugin"
}
```

---

## Contributing

Contributions are welcome. Please open an issue or submit a pull request for enhancements, bug fixes or new plugins. Follow existing code style and add tests where applicable.

---

## License

Forger is released under the MIT License. See [LICENSE](LICENSE) for details.

# Forger

[![Go](https://img.shields.io/badge/go-1.20%2B-blue)](https://go.dev) [![MIT License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Forger is a terminal-native developer dashboard built with Go, [Bubble Tea](https://github.com/charmbracelet/bubbletea), and [Lip Gloss](https://github.com/charmbracelet/lipgloss). It integrates your favorite CLI tools into a unified TUI for static analysis, code navigation, environment inspection, and more.

---

## Features

- **Unified CLI Dashboard**: Seamlessly integrates tools like `marchat`, `ignoregrets`, `CodeSleuth`, `ascii-colorizer`, and more.
- **Plugin Architecture**: Each tool is an independent plugin conforming to a strict interface (`Init`, `Run`, `Render`, `HandleMsg`).
- **Terminal-First UX**: No GUI wrappers or abstractions. Pure, expressive TUI using Bubble Tea and Lip Gloss.
- **Performance-Oriented**: All operations are local and fast. No background network traffic or cloud dependencies.
- **Extensible & Modular**: Add or remove plugins via configuration—no core changes required.
- **Composability**: Plugins can share context (e.g., snapshots or IRs), enabling advanced workflows.
- **Minimalist, Semantic UI**: Crisp layouts, color-rich diagrams (via Mermaid/ascii-colorizer), and contextual overlays.
- **Snapshot & Analysis Tools**: Deep integration of snapshot management (`ignoregrets`), IR visualization (`CodeSleuth`), and chat overlays (`marchat`).

---

## Why Forger?

- **Not an IDE**: Forger is a terminal dashboard—not a replacement for your editor, but a power-tool cockpit for CLI developers.
- **Discoverable**: Clean keybindings, clear plugin panels, and a consistent UX.
- **Fast**: Zero bloat, minimal dependencies, instant startup.

---

## Plugins

- **ignoregrets**: Snapshot manager and viewer.
- **codesleuth**: COBOL IR analyzer and Mermaid-based visualizer.
- **marchat**: Repo-synced terminal chat/log overlay.
- **ascii-colorizer**: GPU-accelerated diagrams and syntax highlighting (planned).
- **parsec**: Log/test output formatter (optional/future).

Each plugin lives in `internal/plugins/` and is loaded via the registry.

---

## Installation

### Prerequisites

- Go 1.20+  
- Terminal supporting true color (for best experience)

### Clone & Build

```sh
git clone https://github.com/Cod-e-Codes/forger
cd forger
go build -o forger ./cmd/forger
```

### Run

```sh
./forger
```

---

## Configuration

Plugins are enabled/configured via `forger.json`:

```json
{
  "default": "ignoregrets",
  "enabled": ["ignoregrets", "codesleuth", "marchat"]
}
```

- Place `forger.json` in the working directory or your config path (`~/.config/forger/`).
- To enable more plugins, add their names to the `enabled` list.
- The `default` field sets the initial active plugin.

---

## Usage & Keybindings

| Key         | Action                                 |
|-------------|----------------------------------------|
| Up/Down     | Navigate plugins in sidebar            |
| c           | Toggle chat overlay (marchat)          |
| Esc         | Close overlay                          |
| q / Ctrl+C  | Quit                                   |

Each plugin may define additional keybindings within its panel.

---

## Architecture

- **Core Runtime**: Manages config, plugin lifecycle, and shared state.
- **Plugin System**: All plugins implement a common Go interface.
- **UI Components**: Shared menus, panels, overlays via Bubble Tea/Lip Gloss.
- **Snapshot Context**: Consistent state for tools like `ignoregrets` and `codesleuth`.

**Directory Structure**
```
cmd/forger/            # Entry point
internal/core/         # Core runtime, context, plugin registry, model
internal/plugins/      # All plugins (marchat, ignoregrets, codesleuth, etc.)
internal/ui/           # Shared UI components (menus, overlays, etc.)
```

---

## Plugin Development

Create a new plugin by implementing the `core.Plugin` interface:

```go
type Plugin interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Plugin, tea.Cmd)
    View() string
    Name() string
}
```
Register your plugin in `internal/core/registry.go`, then add it to your config.

---

## Roadmap

- [x] Core runtime and plugin loader
- [x] Configurable panel layout
- [x] Plugin APIs and lifecycle management
- [x] Integrated snapshot viewer (`ignoregrets`)
- [x] IR visualizer (`codesleuth` + Mermaid)
- [x] Chat overlay (`marchat`)
- [ ] Ascii-colorizer integration
- [ ] Plugin registry & installer
- [ ] Built-in fuzzy finder and command launcher

---

## Contributing

Pull requests and plugin proposals welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for standards and workflow.

---

## License

MIT License © [Cod-e-Codes](https://github.com/Cod-e-Codes)

---

## Acknowledgements

- [Charmbracelet Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Charmbracelet Lip Gloss](https://github.com/charmbracelet/lipgloss)

---

Forger: **Your CLI cockpit. Take control.**

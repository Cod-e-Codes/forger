package main

import (
	"encoding/json"
	"fmt"
	"os"

	"forger/internal/core"

	tea "github.com/charmbracelet/bubbletea"
)

type Config struct {
	Default string   `json:"default"`
	Enabled []string `json:"enabled"`
}

func loadConfig(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func main() {
	cfg, err := loadConfig("forger.json")
	if err != nil {
		core.LogError(fmt.Sprintf("failed to load config: %v", err))
		os.Exit(1)
	}

	model := core.NewModel()
	model.Plugins, model.LoadErrors = core.LoadPlugins(cfg.Enabled, model.Context)

	if _, ok := model.Plugins[cfg.Default]; ok {
		model.Active = cfg.Default
	} else if len(model.Plugins) > 0 {
		model.Active = core.FirstPluginKey(model.Plugins)
	}

	prog := tea.NewProgram(model)
	if _, err := prog.Run(); err != nil {
		core.LogError(fmt.Sprintf("program error: %v", err))
		os.Exit(1)
	}
}

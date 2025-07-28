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

	fmt.Printf("DEBUG: Loading plugins: %v\n", cfg.Enabled)

	model := core.NewModel()
	model.Plugins, model.LoadErrors = core.LoadPlugins(cfg.Enabled, model.Context)

	fmt.Printf("DEBUG: Loaded plugins: %v\n", len(model.Plugins))
	fmt.Printf("DEBUG: Load errors: %v\n", model.LoadErrors)

	if _, ok := model.Plugins[cfg.Default]; ok {
		model.Active = cfg.Default
	} else if len(model.Plugins) > 0 {
		model.Active = core.FirstPluginKey(model.Plugins)
	}

	fmt.Printf("DEBUG: Active plugin: %s\n", model.Active)

	prog := tea.NewProgram(model)
	if _, err := prog.Run(); err != nil {
		core.LogError(fmt.Sprintf("program error: %v", err))
		os.Exit(1)
	}
}

package core

import (
	"fmt"
	"sort"
	"forger/internal/plugins/codesleuth"
	"forger/internal/plugins/ignoregrets"
	"forger/internal/plugins/marchat"
)

// PluginFactory creates a Plugin given shared Context.
type PluginFactory func(ctx *Context) Plugin

// availablePlugins maps plugin names to their factories.
var availablePlugins = map[string]PluginFactory{
	"ignoregrets": ignoregrets.New,
	"codesleuth":  codesleuth.New,
	"marchat":     marchat.New,
	// add ascii-colorizer, parsec, etc.
}

// LoadPlugins instantiates each enabled plugin or records errors.
func LoadPlugins(enabled []string, ctx *Context) (map[string]Plugin, []string) {
	loaded := make(map[string]Plugin)
	var errors []string

	for _, name := range enabled {
		if factory, ok := availablePlugins[name]; ok {
			loaded[name] = factory(ctx)
		} else {
			msg := fmt.Sprintf("plugin '%s' not found in registry", name)
			LogError(msg)
			errors = append(errors, msg)
		}
	}
	return loaded, errors
}

// SortedPluginNames returns plugin names sorted alphabetically.
func SortedPluginNames(plugins map[string]Plugin) []string {
	names := make([]string, 0, len(plugins))
	for name := range plugins {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// FirstPluginKey returns the first key in sorted order.
func FirstPluginKey(plugins map[string]Plugin) string {
	names := SortedPluginNames(plugins)
	if len(names) > 0 {
		return names[0]
	}
	return ""
}

// PrevPluginKey returns the previous plugin key alphabetically.
func PrevPluginKey(plugins map[string]Plugin, current string) string {
	names := SortedPluginNames(plugins)
	for i, name := range names {
		if name == current {
			return names[(i-1+len(names))%len(names)]
		}
	}
	return current
}

// NextPluginKey returns the next plugin key alphabetically.
func NextPluginKey(plugins map[string]Plugin, current string) string {
	names := SortedPluginNames(plugins)
	for i, name := range names {
		if name == current {
			return names[(i+1)%len(names)]
		}
	}
	return current
}

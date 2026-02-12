// Package plugins handles loading and managing the plugin registry.
package plugins

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Plugin represents an available plugin from the registry.
type Plugin struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
}

// Registry holds all available plugins.
type Registry struct {
	Plugins []Plugin `toml:"plugin"`
}

// Load reads the plugin registry from ~/.ocmgr/plugins/plugins.toml.
// Returns an empty registry if the file doesn't exist.
func Load() (*Registry, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home directory: %w", err)
	}

	pluginsFile := filepath.Join(home, ".ocmgr", "plugins", "plugins.toml")

	// Check if file exists
	if _, err := os.Stat(pluginsFile); os.IsNotExist(err) {
		return &Registry{}, nil
	}

	var registry Registry
	if _, err := toml.DecodeFile(pluginsFile, &registry); err != nil {
		return nil, fmt.Errorf("parsing plugins.toml: %w", err)
	}

	return &registry, nil
}

// List returns all available plugins.
func (r *Registry) List() []Plugin {
	return r.Plugins
}

// GetByName finds a plugin by its npm package name.
func (r *Registry) GetByName(name string) *Plugin {
	for i := range r.Plugins {
		if r.Plugins[i].Name == name {
			return &r.Plugins[i]
		}
	}
	return nil
}

// Names returns a slice of all plugin names.
func (r *Registry) Names() []string {
	names := make([]string, len(r.Plugins))
	for i, p := range r.Plugins {
		names[i] = p.Name
	}
	return names
}

// IsEmpty returns true if the registry has no plugins.
func (r *Registry) IsEmpty() bool {
	return len(r.Plugins) == 0
}

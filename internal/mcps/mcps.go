// Package mcps handles loading and managing the MCP server registry.
package mcps

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Definition represents an MCP server definition file.
type Definition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
}

// Registry holds all available MCP servers.
type Registry struct {
	Servers []Definition
}

// Load reads all MCP definitions from ~/.ocmgr/mcps/*.json.
// Returns an empty registry if the directory doesn't exist.
func Load() (*Registry, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home directory: %w", err)
	}

	mcpsDir := filepath.Join(home, ".ocmgr", "mcps")

	// Check if directory exists
	if _, err := os.Stat(mcpsDir); os.IsNotExist(err) {
		return &Registry{}, nil
	}

	// Read all JSON files from the directory
	entries, err := os.ReadDir(mcpsDir)
	if err != nil {
		return nil, fmt.Errorf("reading mcps directory: %w", err)
	}

	var servers []Definition
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(mcpsDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", entry.Name(), err)
		}

		var def Definition
		if err := json.Unmarshal(data, &def); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}

		servers = append(servers, def)
	}

	return &Registry{Servers: servers}, nil
}

// List returns all available MCP servers.
func (r *Registry) List() []Definition {
	return r.Servers
}

// GetByName finds an MCP by its name field.
func (r *Registry) GetByName(name string) *Definition {
	for i := range r.Servers {
		if r.Servers[i].Name == name {
			return &r.Servers[i]
		}
	}
	return nil
}

// Names returns a slice of all MCP server names.
func (r *Registry) Names() []string {
	names := make([]string, len(r.Servers))
	for i, s := range r.Servers {
		names[i] = s.Name
	}
	return names
}

// IsEmpty returns true if the registry has no MCP servers.
func (r *Registry) IsEmpty() bool {
	return len(r.Servers) == 0
}

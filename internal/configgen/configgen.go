// Package configgen handles generating opencode.json configuration files.
package configgen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the opencode.json structure.
type Config struct {
	Schema string              `json:"$schema,omitempty"`
	Plugin []string            `json:"plugin,omitempty"`
	MCP    map[string]MCPEntry `json:"mcp,omitempty"`
}

// MCPEntry represents an MCP server entry in opencode.json.
type MCPEntry struct {
	Type        string            `json:"type"`
	Command     []string          `json:"command,omitempty"`
	URL         string            `json:"url,omitempty"`
	Enabled     bool              `json:"enabled,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	OAuth       interface{}       `json:"oauth,omitempty"`
	Timeout     int               `json:"timeout,omitempty"`
}

// Options for generating the config file.
type Options struct {
	// Plugins to include in the config.
	Plugins []string
	// MCPs to include, keyed by name.
	MCPs map[string]MCPEntry
}

// NewConfig creates a new Config with the schema already set.
func NewConfig() *Config {
	return &Config{
		Schema: "https://opencode.ai/config.json",
	}
}

// AddPlugin adds a plugin to the config.
func (c *Config) AddPlugin(name string) {
	c.Plugin = append(c.Plugin, name)
}

// AddMCP adds an MCP server to the config.
func (c *Config) AddMCP(name string, entry MCPEntry) {
	if c.MCP == nil {
		c.MCP = make(map[string]MCPEntry)
	}
	c.MCP[name] = entry
}

// HasPlugins returns true if there are any plugins configured.
func (c *Config) HasPlugins() bool {
	return len(c.Plugin) > 0
}

// HasMCPs returns true if there are any MCPs configured.
func (c *Config) HasMCPs() bool {
	return len(c.MCP) > 0
}

// IsEmpty returns true if there are no plugins or MCPs.
func (c *Config) IsEmpty() bool {
	return !c.HasPlugins() && !c.HasMCPs()
}

// Write writes the config to the specified directory as opencode.json.
func (c *Config) Write(targetDir string) error {
	// Ensure directory exists
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	filePath := filepath.Join(targetDir, "opencode.json")

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	data = append(data, '\n')

	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

// Load reads an existing opencode.json from the specified directory.
// Returns nil if the file doesn't exist.
func Load(targetDir string) (*Config, error) {
	filePath := filepath.Join(targetDir, "opencode.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &config, nil
}

// Merge combines another config into this one.
// Plugins are deduplicated, MCPs are merged (existing keys are preserved).
func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}

	// Merge plugins (deduplicate)
	existingPlugins := make(map[string]bool)
	for _, p := range c.Plugin {
		existingPlugins[p] = true
	}
	for _, p := range other.Plugin {
		if !existingPlugins[p] {
			c.Plugin = append(c.Plugin, p)
		}
	}

	// Merge MCPs (existing keys take precedence)
	if other.MCP != nil {
		if c.MCP == nil {
			c.MCP = make(map[string]MCPEntry)
		}
		for name, entry := range other.MCP {
			if _, exists := c.MCP[name]; !exists {
				c.MCP[name] = entry
			}
		}
	}
}

// Generate creates an opencode.json file with the specified options.
// If a file already exists, it merges the new config with the existing one.
func Generate(targetDir string, opts Options) error {
	// Load existing config if it exists
	config, err := Load(targetDir)
	if err != nil {
		return err
	}

	if config == nil {
		config = NewConfig()
	}

	// Add plugins
	for _, p := range opts.Plugins {
		config.AddPlugin(p)
	}

	// Add MCPs
	for name, entry := range opts.MCPs {
		config.AddMCP(name, entry)
	}

	// Only write if there's something to write
	if config.IsEmpty() {
		return nil
	}

	return config.Write(targetDir)
}

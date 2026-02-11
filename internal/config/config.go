// Package config manages the global ocmgr configuration file (~/.ocmgr/config.toml).
package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config is the top-level configuration for ocmgr.
type Config struct {
	GitHub   GitHub   `toml:"github"`
	Defaults Defaults `toml:"defaults"`
	Store    Store    `toml:"store"`
}

// GitHub holds settings for the remote profile repository.
type GitHub struct {
	// Repo is the owner/repo slug on GitHub (e.g. "acchapm1/opencode-profiles").
	Repo string `toml:"repo"`
	// Auth is the authentication method: "gh", "env", "ssh", or "token".
	Auth string `toml:"auth"`
}

// Defaults holds user-facing default behaviours.
type Defaults struct {
	// MergeStrategy controls how conflicting files are handled.
	// One of "prompt", "overwrite", "merge", or "skip".
	MergeStrategy string `toml:"merge_strategy"`
	// Editor is the command used to open files for editing.
	Editor string `toml:"editor"`
}

// Store holds settings for the local profile store.
type Store struct {
	// Path is the directory where downloaded profiles are kept.
	// The "~" prefix is expanded to the user's home directory at runtime.
	Path string `toml:"path"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		GitHub: GitHub{
			Repo: "acchapm1/opencode-profiles",
			Auth: "gh",
		},
		Defaults: Defaults{
			MergeStrategy: "prompt",
			Editor:        "nvim",
		},
		Store: Store{
			Path: "~/.ocmgr/profiles",
		},
	}
}

// ConfigDir returns the absolute path to the ocmgr configuration directory
// (~/.ocmgr with the tilde expanded).
func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fall back to the HOME env var; if that is also empty the caller
		// will get a relative path, which is the best we can do.
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".ocmgr")
}

// ConfigPath returns the absolute path to the ocmgr configuration file
// (~/.ocmgr/config.toml).
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.toml")
}

// Load reads the configuration from disk. If the file does not exist the
// default configuration is returned without an error.
func Load() (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if _, err := toml.Decode(string(data), cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save writes cfg to ~/.ocmgr/config.toml, creating the configuration
// directory if it does not already exist.
func Save(cfg *Config) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(cfg); err != nil {
		return err
	}

	return os.WriteFile(ConfigPath(), buf.Bytes(), 0o644)
}

// EnsureConfigDir creates the ~/.ocmgr directory (and any parents) if it does
// not already exist.
func EnsureConfigDir() error {
	return os.MkdirAll(ConfigDir(), 0o755)
}

// ExpandPath replaces a leading "~/" or bare "~" in path with the current
// user's home directory. Paths like "~user" are not expanded and are returned
// unchanged.
func ExpandPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}

	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	return path
}

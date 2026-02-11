// Package profile defines the Profile data model and operations for
// reading, writing, validating, scaffolding, and listing profile contents.
//
// A profile is a directory (e.g. ~/.ocmgr/profiles/go/) that contains a
// profile.toml metadata file and up to four content subdirectories:
// agents/, commands/, skills/, and plugins/.
package profile

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

// validName matches profile names that are safe directory names:
// alphanumeric, hyphens, underscores, and dots, starting with an alphanumeric.
var validName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)

// ValidateName checks that a profile name is safe to use as a directory name.
// It rejects empty names, path traversal attempts, and special characters.
func ValidateName(name string) error {
	if name == "" {
		return errors.New("profile name must not be empty")
	}
	if name == "." || name == ".." || strings.ContainsAny(name, "/\\") || strings.Contains(name, "..") {
		return fmt.Errorf("invalid profile name %q: must be a simple directory name", name)
	}
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid profile name %q: must start with alphanumeric and contain only alphanumeric, hyphens, underscores, or dots", name)
	}
	return nil
}

// Profile represents the metadata and location of an ocmgr profile.
type Profile struct {
	// Name is the short identifier for the profile (required).
	Name string `toml:"name"`
	// Description is a human-readable summary of the profile.
	Description string `toml:"description"`
	// Version is a semver-style version string.
	Version string `toml:"version"`
	// Author is the profile creator's identifier.
	Author string `toml:"author"`
	// Tags is an optional list of keywords for discovery.
	Tags []string `toml:"tags"`
	// Extends names another profile that this one inherits from.
	Extends string `toml:"extends"`
	// Path is the absolute directory path on disk. It is not serialized to TOML.
	Path string `toml:"-"`
}

// profileTOML is the on-disk TOML representation that wraps Profile
// in a [profile] table.
type profileTOML struct {
	Profile Profile `toml:"profile"`
}

// Contents describes the files found inside a profile's content directories.
type Contents struct {
	// Agents lists relative paths to *.md files under agents/.
	Agents []string
	// Commands lists relative paths to *.md files under commands/.
	Commands []string
	// Skills lists relative paths to SKILL.md files under skills/<name>/.
	Skills []string
	// Plugins lists relative paths to *.ts files under plugins/.
	Plugins []string
	// HasPackageJSON indicates whether plugins/package.json exists.
	HasPackageJSON bool
}

// ContentDirs returns the four content subdirectory names that a profile
// may contain.
func ContentDirs() []string {
	return []string{"agents", "commands", "skills", "plugins"}
}

// LoadProfile reads profile.toml from dir and returns the parsed Profile.
// The returned Profile's Path field is set to the absolute path of dir.
func LoadProfile(dir string) (*Profile, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolving profile directory: %w", err)
	}

	tomlPath := filepath.Join(absDir, "profile.toml")

	data, err := os.ReadFile(tomlPath)
	if err != nil {
		return nil, fmt.Errorf("reading profile.toml: %w", err)
	}

	var doc profileTOML
	if err := toml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing profile.toml: %w", err)
	}

	p := &doc.Profile
	p.Path = absDir
	return p, nil
}

// SaveProfile writes p to profile.toml inside p.Path, creating the
// directory (and parents) if it does not already exist.
func SaveProfile(p *Profile) error {
	if p.Path == "" {
		return errors.New("profile path is empty")
	}

	if err := os.MkdirAll(p.Path, 0o755); err != nil {
		return fmt.Errorf("creating profile directory: %w", err)
	}

	doc := profileTOML{Profile: *p}

	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(doc); err != nil {
		return fmt.Errorf("encoding profile.toml: %w", err)
	}

	tomlPath := filepath.Join(p.Path, "profile.toml")
	if err := os.WriteFile(tomlPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("writing profile.toml: %w", err)
	}

	return nil
}

// Validate checks that a profile is well-formed:
//   - Name must be non-empty.
//   - Path must exist on disk and be a directory.
//   - At least one content directory (agents/, commands/, skills/, plugins/)
//     must exist and contain at least one entry.
func Validate(p *Profile) error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("profile name must not be empty")
	}

	info, err := os.Stat(p.Path)
	if err != nil {
		return fmt.Errorf("profile path %q: %w", p.Path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("profile path %q is not a directory", p.Path)
	}

	hasContent := false
	for _, d := range ContentDirs() {
		dirPath := filepath.Join(p.Path, d)
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			// Directory does not exist or is unreadable — skip.
			continue
		}
		if len(entries) > 0 {
			hasContent = true
			break
		}
	}

	if !hasContent {
		return fmt.Errorf("profile %q has no content: at least one of %v must exist and be non-empty",
			p.Name, ContentDirs())
	}

	return nil
}

// ListContents scans the profile directory and returns a Contents struct
// describing every content file found. Paths in the returned slices are
// relative to the profile root (e.g. "agents/code-reviewer.md").
func ListContents(p *Profile) (*Contents, error) {
	c := &Contents{}

	// agents/ — top-level *.md files.
	agents, err := listMD(filepath.Join(p.Path, "agents"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("listing agents: %w", err)
	}
	for _, name := range agents {
		c.Agents = append(c.Agents, filepath.Join("agents", name))
	}

	// commands/ — top-level *.md files.
	commands, err := listMD(filepath.Join(p.Path, "commands"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("listing commands: %w", err)
	}
	for _, name := range commands {
		c.Commands = append(c.Commands, filepath.Join("commands", name))
	}

	// skills/ — each subdirectory contains a SKILL.md.
	skillsDir := filepath.Join(p.Path, "skills")
	skillEntries, err := os.ReadDir(skillsDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("listing skills: %w", err)
	}
	for _, entry := range skillEntries {
		if !entry.IsDir() {
			continue
		}
		skillFile := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
		if _, err := os.Stat(skillFile); err == nil {
			c.Skills = append(c.Skills, filepath.Join("skills", entry.Name(), "SKILL.md"))
		}
	}

	// plugins/ — top-level *.ts files and optional package.json.
	pluginsDir := filepath.Join(p.Path, "plugins")
	pluginEntries, err := os.ReadDir(pluginsDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("listing plugins: %w", err)
	}
	for _, entry := range pluginEntries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if name == "package.json" {
			c.HasPackageJSON = true
			continue
		}
		if strings.HasSuffix(name, ".ts") {
			c.Plugins = append(c.Plugins, filepath.Join("plugins", name))
		}
	}

	return c, nil
}

// ScaffoldProfile creates an empty profile directory at dir/<name>
// containing a profile.toml and the four empty content subdirectories.
// It returns the newly created Profile.
func ScaffoldProfile(dir string, name string) (*Profile, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolving scaffold directory: %w", err)
	}

	profileDir := filepath.Join(absDir, name)

	// Create content subdirectories (this implicitly creates profileDir).
	for _, sub := range ContentDirs() {
		if err := os.MkdirAll(filepath.Join(profileDir, sub), 0o755); err != nil {
			return nil, fmt.Errorf("creating %s directory: %w", sub, err)
		}
	}

	p := &Profile{
		Name: name,
		Path: profileDir,
	}

	if err := SaveProfile(p); err != nil {
		return nil, fmt.Errorf("writing scaffold profile.toml: %w", err)
	}

	return p, nil
}

// listMD returns the names of all *.md files in the given directory.
// If the directory does not exist the underlying os error is returned
// so callers can check with errors.Is(err, os.ErrNotExist).
func listMD(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".md") {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

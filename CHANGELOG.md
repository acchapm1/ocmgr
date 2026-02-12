# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [Unreleased]

### Added

- **Plugin selection feature plan** (`plugintodo.md`)
  - Phase 1: Remove opencode.json from copy process
  - Phase 2: Interactive plugin selection during `ocmgr init`
  - Plugins loaded from `~/.ocmgr/plugins/plugins.toml`

- **MCP server selection feature plan** (`mcptodo.md`)
  - Phase 3: Interactive MCP server selection during `ocmgr init`
  - MCPs loaded from individual JSON files in `~/.ocmgr/mcps/`
  - Supports npx-based, Docker-based, and remote MCP servers

- **Sample MCP definition files** in `~/.ocmgr/mcps/`
  - `sequentialthinking.json` - npx-based sequential thinking MCP
  - `context7.json` - npx-based documentation search MCP
  - `docker.json` - Docker-based container management MCP
  - `sentry.json` - Remote Sentry integration MCP
  - `filesystem.json` - npx-based filesystem operations MCP

- **Plugin registry file** (`~/.ocmgr/plugins/plugins.toml`)
  - TOML format with name and description for each plugin
  - Replaces plain text `plugins` file

### Changed

- **USAGE.md** - Comprehensive documentation updates
  - Added Table of Contents entries for all commands
  - Added `ocmgr profile import` command documentation
  - Added `ocmgr profile export` command documentation
  - Added `ocmgr sync push` command documentation
  - Added `ocmgr sync pull` command documentation
  - Added `ocmgr sync status` command documentation
  - Updated "Sharing Profiles" workflow section with full GitHub sync instructions
  - Removed "(Phase 2)" references from config documentation

- **internal/copier/copier.go** - Added support for root-level profile files
  - Added `profileFiles` map for recognised root-level files (e.g., `opencode.json`)
  - Updated walk filter logic to allow root-level files through to copy logic
  - Updated include/exclude filtering to not block root-level profile files
  - Updated doc comments to reflect new behavior
  - `opencode.json` is now copied during `ocmgr init` (to be removed in Phase 1 of plugin feature)

---

## [0.1.0] - 2025-02-12

### Added

- Initial release of ocmgr
- Core profile management commands:
  - `ocmgr init` - Initialize .opencode/ from profile(s)
  - `ocmgr profile list` - List all local profiles
  - `ocmgr profile show` - Show profile details
  - `ocmgr profile create` - Create new empty profile
  - `ocmgr profile delete` - Delete a profile
  - `ocmgr profile import` - Import from local dir or GitHub URL
  - `ocmgr profile export` - Export to a directory
  - `ocmgr snapshot` - Capture .opencode/ as a new profile
  - `ocmgr config show` - Show current configuration
  - `ocmgr config set` - Set a config value
  - `ocmgr config init` - Interactive first-run setup

- GitHub sync commands:
  - `ocmgr sync push` - Push profile to GitHub
  - `ocmgr sync pull` - Pull profile(s) from GitHub
  - `ocmgr sync status` - Show sync status

- Profile features:
  - Profile inheritance via `extends` in profile.toml
  - Selective init with `--only` and `--exclude` flags
  - Conflict resolution: interactive prompt, `--force`, `--merge`
  - Dry-run mode with `--dry-run`
  - Plugin dependency detection (prompts for `bun install`)

- Installation methods:
  - curl installer (`install.sh`)
  - Homebrew formula (`Formula/ocmgr.rb`)
  - Build from source via Makefile

- GitHub Actions workflow for release binaries (`.github/workflows/release.yml`)

---

## Future Roadmap

### Phase 1 - Remove opencode.json Copy
- Remove `opencode.json` from `profileFiles` map
- Delete template `opencode.json` files from profiles

### Phase 2 - Plugin Selection
- Interactive plugin selection during `ocmgr init`
- Generate `opencode.json` with selected plugins

### Phase 3 - MCP Server Selection
- Interactive MCP server selection during `ocmgr init`
- Merge MCP configs into generated `opencode.json`

### Phase 4 - Interactive TUI
- Full TUI built with Charmbracelet (Bubble Tea, Huh, Lip Gloss)
- Profile browser, init wizard, diff viewer, snapshot wizard

### Phase 5 - Advanced Features
- Profile registry
- Template variables
- Pre/post init hooks
- Auto-detect project type
- Shell completions
- `ocmgr doctor` command

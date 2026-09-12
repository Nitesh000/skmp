# Changelog

All notable changes to skmp will be documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Versioning follows [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

### Added

- Per-skill harness access control: toggle individual harnesses on/off per skill from the detail pane (`tab` to focus, `j`/`k` to navigate, `space`/`enter` to toggle)
- Mouse support throughout the TUI: scroll to navigate lists and detail pane, click tabs to switch, click list items to select, click harness rows to toggle access
- `ctrl+o` key — open selected skill or bundle in browser (converts raw source URL to viewable GitHub tree URL)
- `R` key — report skill issue, opens a prefilled GitHub issue with skill name, version, author, and source pre-populated
- `agy` harness support (Antigravity CLI) — skills path `~/.gemini/config/skills/`
- Three new harness functions: `SkillHarnessState`, `LinkSkillTo`, `UnlinkSkillFrom`
- Help menu updated with all new keybindings

### Fixed

- Cursor skills path corrected from `~/.cursor/skills/` to `~/.cursor/skills-cursor/`
- Cursor installed detection changed from missing `mcp.json` check to `cursor` in PATH
- Detail pane overflow: `scroll()` now always truncates to pane height instead of passing full content through when offset is zero
- Antigravity harness renamed to `antigravity-ide` to distinguish from `agy` CLI

---

## [0.1.0] — 2025-09-11

### Added

- Full-screen TUI with Skills, Bundles, My Skills, and My Bundles tabs
- `skmp add` — install one or more skills by name
- `skmp add --bundle` — install all skills in a named bundle
- `skmp remove` — uninstall one or more skills
- `skmp list` — list installed skills
- `skmp sync` — re-link skills into newly installed harnesses
- `skmp update` — force-refresh the local registry cache
- `skmp version` — print binary version
- Harness support: pi, Claude Code, Antigravity IDE, Codex, Cursor, OpenCode
- Correct per-harness skills paths:
  - pi → `~/.agents/skills/`
  - Claude Code → `~/.claude/skills/`
  - Antigravity → `~/.gemini/antigravity-ide/skills/`
  - Codex → `~/.codex/skills/`
  - Cursor → `~/.cursor/skills/`
  - OpenCode → registered in `~/.config/opencode/opencode.jsonc`
- Registry backed by `registry/index.json`, served directly from GitHub
- Local 24 h registry cache
- Bleve full-text search with wildcard + fuzzy queries
- Custom bleve analyzer that preserves hyphens, dots, and other punctuation in skill names
- Symlink-based install (copy fallback on Windows)
- `MaxHeight` on TUI panes — tab bar, status bar, and help line always visible
- List pane auto-scrolls to keep selected item in view
- Seed skills: caveman, diagnose, tdd, write-a-skill
- Seed bundle: nitesh000/core
- NPM distribution (`npm install -g skmp`)
- Homebrew distribution via GoReleaser (tap: `Nitesh000/homebrew-tap`)

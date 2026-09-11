# Changelog

All notable changes to skmp will be documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Versioning follows [Semantic Versioning](https://semver.org/).

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
- Registry backed by `registry/index.json`, served via jsDelivr CDN
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

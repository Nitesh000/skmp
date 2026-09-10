# Changelog

All notable changes to skmp will be documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Versioning follows [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

### Added
- Full-screen TUI with skills and bundles tabs
- `skmp add` — install one or more skills
- `skmp remove` — uninstall one or more skills
- `skmp list` — list installed skills
- `skmp sync` — re-link skills into newly installed harnesses
- `--bundle` flag on `skmp add` — install all skills in a bundle
- Harness support: pi, Claude Code, OpenCode, Codex, Antigravity
- Registry backed by `registry/index.json`, served via jsDelivr CDN
- Local 24h cache for registry index
- Fuzzy search via Bleve
- Symlink-based install (copy fallback on Windows)
- Seed skills: caveman, diagnose, tdd, write-a-skill
- Seed bundle: nitesh000/core

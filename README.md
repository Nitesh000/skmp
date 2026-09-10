# skmp

A full-screen TUI package manager for AI agent skills.

Install, search, and manage skills across all major agentic tools — pi, Claude Code, OpenCode, Codex, and Antigravity — from one place.

```
┌─ skmp ──────────────────────────────────────────────────────────┐
│ [1] Skills  [2] Bundles          / search...        [?] help    │
├──────────────────────┬──────────────────────────────────────────┤
│ Skills (4)           │ caveman                                  │
│                      │                                          │
│ ▶ [✓] caveman        │ Ultra-compressed communication mode.     │
│   [✓] diagnose       │ Cuts token usage ~75% by dropping        │
│   [ ] tdd            │ filler and pleasantries.                 │
│   [ ] write-a-skill  │                                          │
│                      │ Version   1.0.0                          │
│                      │ Author    nitesh000                      │
│                      │ Tags      communication  tokens  terse   │
│                      │                                          │
│                      │ ● installed                              │
│                      │ [x] remove                               │
├──────────────────────┴──────────────────────────────────────────┤
│ Installed: 2 │ Available: 4          ✓ ready      skmp v0.0.1   │
└─────────────────────────────────────────────────────────────────┘
```

## Install

```bash
npm install -g skmp
```

Or download a binary directly from [Releases](https://github.com/nitesh000/skill-set/releases).

```bash
# macOS (Homebrew) — coming soon
brew install skmp
```

## Usage

### TUI (recommended)

```bash
skmp
```

Opens the full-screen interface. Navigate with keyboard, install/remove with a keypress.

### CLI

```bash
# install skills
skmp add caveman
skmp add caveman diagnose tdd       # multiple at once
skmp add --bundle nitesh000/core    # install a whole bundle

# remove skills
skmp remove caveman
skmp remove caveman diagnose

# list installed skills
skmp list

# re-link skills into a newly installed harness
skmp sync
```

## Keybindings

| Key         | Action            |
| ----------- | ----------------- |
| `j` / `↓`   | move down         |
| `k` / `↑`   | move up           |
| `1`         | skills tab        |
| `2`         | bundles tab       |
| `/`         | search            |
| `esc`       | clear search      |
| `tab` / `l` | focus detail pane |
| `h`         | focus list pane   |
| `i`         | install selected  |
| `x`         | remove selected   |
| `?`         | toggle help       |
| `q`         | quit              |

## Supported Harnesses

| Harness                                  | Skills path         |
| ---------------------------------------- | ------------------- |
| [pi](https://github.com/pi-cli/pi)       | `~/.agents/skills/` |
| [Claude Code](https://claude.ai/code)    | `~/.agents/skills/` |
| [OpenCode](https://opencode.ai)          | `~/.agents/skills/` |
| [Antigravity](https://antigravity.dev)   | `~/.agents/skills/` |
| [Codex](https://github.com/openai/codex) | `~/.codex/skills/`  |

Skills are installed once to `~/.skmp/skills/` and symlinked into each harness automatically.
Run `skmp sync` after installing a new harness to link existing skills into it.

## Registry

The skill registry lives at [`registry/index.json`](registry/index.json) in this repo and is served via jsDelivr CDN. It is cached locally for 24 hours.

Skills themselves live in their **author's own GitHub repo** — the registry only stores metadata and a source URL. skmp fetches skill files on demand at install time.

### Submitting a skill

1. Create a GitHub repo with at least a `SKILL.md` following the [skill format](#skill-format)
2. Open a PR adding your skill to `registry/index.json`
3. That's it — you host the files, the registry just indexes them

### Submitting a bundle

A bundle is a named collection of skills from one repo. Add a `bundles` entry to `registry/index.json` pointing to your repo:

```json
{
  "name": "yourname/bundle-name",
  "description": "What this bundle does",
  "author": "yourname",
  "repo": "https://github.com/yourname/your-repo",
  "branch": "main",
  "skills_path": "skills",
  "skills": ["skill-one", "skill-two"]
}
```

## Skill Format

Every skill needs at minimum a `SKILL.md`:

```markdown
---
name: your-skill
description: One line description of what the skill does.
---

Instructions for the AI agent go here.
```

Optional files: `REFERENCE.md`, `EXAMPLES.md`

## How It Works

```
registry/index.json   — metadata + source URLs (hosted here, served via CDN)
~/.skmp/skills/       — downloaded skill files (skmp owns this)
~/.agents/skills/     — symlinks → ~/.skmp/skills/ (read by most harnesses)
~/.codex/skills/      — symlinks → ~/.skmp/skills/ (read by Codex)
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE).

# Contributing to skmp

There are three ways to contribute:

1. **Add a skill** — publish a skill to the registry
2. **Add a bundle** — group existing skills into a named collection
3. **Improve skmp itself** — fix bugs, add features to the CLI/TUI

---

## Adding a Skill

### 1. Create your skill

Host your skill files in your own GitHub repo. Minimum required file: `SKILL.md`.

```
your-repo/
└── your-skill/
    ├── SKILL.md        ← required
    ├── REFERENCE.md    ← optional
    └── EXAMPLES.md     ← optional
```

`SKILL.md` format:

```markdown
---
name: your-skill
description: One line description. Used in search results.
---

Instructions for the AI agent go here.
Explain when to use this skill and how it behaves.
```

The `source` URL you'll submit should point to the folder containing `SKILL.md`:
```
https://raw.githubusercontent.com/yourname/your-repo/main/your-skill/
```

### 2. Add to the registry

Open `registry/index.json` and add your skill to the `skills` array:

```json
{
  "name": "your-skill",
  "description": "One line description matching your SKILL.md frontmatter.",
  "version": "1.0.0",
  "author": "your-github-username",
  "tags": ["tag1", "tag2"],
  "harnesses": ["all"],
  "source": "https://raw.githubusercontent.com/yourname/your-repo/main/your-skill/"
}
```

**Fields:**
- `name` — kebab-case, unique across the registry
- `description` — shown in search results and TUI detail pane
- `version` — semver, matches what's in your repo
- `author` — your GitHub username
- `tags` — 2–5 tags that describe the skill's purpose
- `harnesses` — `["all"]` works for most skills; specify `["pi", "claude-code"]` etc. if it's harness-specific
- `source` — raw GitHub URL to the folder containing `SKILL.md` (trailing slash required)

### 3. Open a PR

```bash
git clone https://github.com/Nitesh000/skmp
cd skmp
# edit registry/index.json
git checkout -b add-skill-your-skill-name
git add registry/index.json
git commit -m "feat: add skill your-skill-name"
git push origin add-skill-your-skill-name
```

PR checklist (verified before merge):

- [ ] `SKILL.md` is reachable at the `source` URL
- [ ] `name` is unique in the registry
- [ ] `description` is one sentence, clear, and matches `SKILL.md`
- [ ] `source` URL ends with `/`
- [ ] Skill works with at least one harness

---

## Adding a Bundle

A bundle groups skills from one repo into a named collection.

Add to the `bundles` array in `registry/index.json`:

```json
{
  "name": "yourname/bundle-name",
  "description": "Short description of what this bundle is for.",
  "author": "yourname",
  "repo": "https://github.com/yourname/your-repo",
  "branch": "main",
  "skills_path": "path/to/skills",
  "skills": ["skill-one", "skill-two", "skill-three"]
}
```

**Fields:**
- `name` — `author/bundle-name` format, must be unique
- `skills_path` — path inside your repo where skill folders live (e.g. `skills`, `registry/skills`)
- `skills` — list of skill names in the bundle; each will be fetched from `repo/branch/skills_path/<name>/`

Every skill listed in a bundle **must also be individually listed** in the `skills` array. Bundles are a convenience shortcut — every skill must be independently installable via `skmp add <name>`. PRs that add a bundle without corresponding individual skill entries will be asked to add them before merge.

---

## Contributing to skmp (the tool)

### Setup

```bash
git clone https://github.com/Nitesh000/skmp
cd skmp/skmp
go mod download
go build ./...
```

### Running locally

```bash
go run main.go           # opens TUI
go run main.go list      # list command
go run main.go add tdd   # add command
```

### Project structure

```
skmp/
├── cmd/          — Cobra CLI commands (one file per command)
├── config/       — version constant
├── registry/     — index fetching, caching, search (Bleve)
├── harness/      — harness detection, skill install/remove/sync
└── tui/          — Bubbletea TUI (app.go = model/update/view, styles.go = lipgloss styles)
```

### Supported harnesses

| Harness | Detection | Skills path |
| --- | --- | --- |
| pi | `pi` in PATH | `~/.agents/skills/` |
| Claude Code | `claude` in PATH | `~/.claude/skills/` |
| Antigravity | `~/.antigravity-ide/antigravity-ide/bin/agy-ide` exists | `~/.gemini/antigravity-ide/skills/` |
| Codex | `codex` in PATH | `~/.codex/skills/` |
| Cursor | `~/.cursor/mcp.json` exists | `~/.cursor/skills/` |
| OpenCode | `opencode` in PATH | config-based (`opencode.jsonc`) |

### Guidelines

- Keep commands in `cmd/` thin — logic belongs in `registry/` or `harness/`
- No new dependencies without discussion — binary size matters
- Cross-platform: be mindful of Windows paths and symlinks (copy fallback exists)
- Error messages should tell the user what to do next, not just what went wrong

### Commit style

```
feat: add antigravity harness support
fix: correct claude-code skills path to ~/.claude/skills
docs: update supported harnesses table
refactor: deduplicate harness symlink logic
```

### Opening a PR

- One PR per feature or fix
- Include a short description of what changed and why
- If it's a registry change (new skill/bundle), see the sections above

---

## Questions

Open an issue or start a discussion on GitHub.

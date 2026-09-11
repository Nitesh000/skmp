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
- `tags` — 2-5 tags that describe the skill's purpose
- `harnesses` — `["all"]` works for most skills; specify `["pi", "claude-code"]` etc. if it's harness-specific
- `source` — raw GitHub URL to the folder containing `SKILL.md` (trailing slash required)

### 3. Open a PR

```bash
git clone https://github.com/nitesh000/skill-set
cd skill-set
# edit registry/index.json
git checkout -b add-skill-your-skill-name
git add registry/index.json
git commit -m "feat: add skill your-skill-name"
git push origin add-skill-your-skill-name
```

Open a PR. The checklist will be verified before merge:

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
- `skills` — list of skill names in the bundle; each will be downloaded from `repo/branch/skills_path/<name>/`

Every skill listed in a bundle **must also be individually listed** in the `skills` array. Bundles are a convenience shortcut, not a namespace — every skill must be independently installable via `skmp add <name>`. PRs that add a bundle without corresponding individual skill entries will be asked to add them before merge.

---

## Contributing to skmp (the tool)

### Setup

```bash
git clone https://github.com/nitesh000/skill-set
cd skill-set/skmp
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
├── registry/     — index fetching, caching, search
├── harness/      — harness detection, skill install/remove/sync
└── tui/          — Bubbletea TUI (app.go = model/update/view, styles.go = lipgloss styles)
```

### Guidelines

- Keep commands in `cmd/` thin — logic belongs in `registry/` or `harness/`
- No new dependencies without discussion — the binary size matters
- Cross-platform: test on Mac + Linux, be mindful of Windows paths and symlinks
- Error messages should tell the user what to do next, not just what went wrong

### Commit style

```
feat: add skmp sync command
fix: handle missing ~/.agents/skills dir gracefully
docs: update contributing guide
refactor: deduplicate harness symlink logic
```

### Opening a PR

- One PR per feature or fix
- Include a short description of what changed and why
- If it's a registry change (new skill/bundle), see the sections above

---

## Questions

Open an issue or start a discussion on GitHub.

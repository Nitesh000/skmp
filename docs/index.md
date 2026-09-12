---
layout: default
title: Home
nav_order: 1
description: "skmp — A TUI package manager for AI agent skills"
permalink: /
---

# skmp
{: .fs-9 }

A TUI package manager for AI agent skills. Install, search, and manage skills across all major agentic tools from one place.
{: .fs-6 .fw-300 }

[Get Started]({{ site.baseurl }}/getting-started/){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View on GitHub](https://github.com/Nitesh000/skmp){: .btn .fs-5 .mb-4 .mb-md-0 }
[Download]({{ site.baseurl }}/download/){: .btn .fs-5 .mb-4 .mb-md-0 }

---

```bash
npm install -g skmp-cli
```

## Supported agents

pi · Claude Code · Antigravity IDE · agy · Codex · Cursor · OpenCode
{: .fs-5 .fw-300 }

---

![Skills tab]({{ site.baseurl }}/images/skills.png)
*Browse, search, and install skills with a single keypress*

---

## Features

### Full-text search
Fuzzy + wildcard search powered by Bleve. Press `/` to find any skill instantly.

### Auto-linking
Skills are downloaded once to `~/.skmp/skills/` and symlinked into every detected harness. No duplicate files, no manual config.

### Per-skill harness control
Toggle which agents have access to each skill. Focus the detail pane with `tab`, navigate with `j`/`k`, toggle with `space`.

### Bundles
Install curated collections of skills with one command: `skmp add --bundle author/bundle`

### Mouse + keyboard
Full mouse support — scroll, click tabs, click list items, click harness toggles. Vim-style keybindings throughout.

### Open registry
Skills live in their author's repo. The registry just indexes metadata and source URLs. Submit a PR to add yours.

---

![Bundles tab]({{ site.baseurl }}/images/bundle.png)
*Browse bundles — curated skill collections*

---

## How it works

```
registry/index.json                    — metadata + source URLs (CDN-served)
~/.skmp/skills/                        — downloaded skill files (skmp owns this)
~/.agents/skills/                      — symlinks → ~/.skmp/skills/  (pi)
~/.claude/skills/                      — symlinks → ~/.skmp/skills/  (Claude Code)
~/.gemini/antigravity-ide/skills/      — symlinks → ~/.skmp/skills/  (Antigravity IDE)
~/.gemini/config/skills/               — symlinks → ~/.skmp/skills/  (agy)
~/.codex/skills/                       — symlinks → ~/.skmp/skills/  (Codex)
~/.cursor/skills-cursor/               — symlinks → ~/.skmp/skills/  (Cursor)
~/.config/opencode/opencode.jsonc      — skills.paths entry (OpenCode)
```

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `↓` | move down (list) / next harness (detail) |
| `k` / `↑` | move up (list) / prev harness (detail) |
| `K` / `J` | jump to top / bottom |
| `1`–`4` | switch tabs |
| `tab` | toggle list / detail focus |
| `space` / `enter` | toggle harness access (detail focused) |
| `/` | search |
| `i` / `x` | install / remove |
| `ctrl+o` | open in browser |
| `R` | report skill issue |
| `?` | help |
| `q` | quit |

Mouse is fully supported: scroll to navigate, click tabs to switch, click list items to select, click harness rows to toggle access.

---

![Help menu]({{ site.baseurl }}/images/help.png)
*Press `?` for the full help overlay*

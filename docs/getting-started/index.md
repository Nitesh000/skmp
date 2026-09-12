---
layout: default
title: Getting Started
nav_order: 2
description: Install skmp and manage AI agent skills in under a minute.
---

# Getting Started
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Install

### npm (recommended)

```bash
npm install -g skmp-cli
```

### Homebrew (macOS / Linux)

```bash
brew tap Nitesh000/tap
brew install skmp
```

### Binary download

```bash
curl -L https://github.com/Nitesh000/skmp/releases/latest/download/skmp_$(uname -s)_$(uname -m).tar.gz | tar xz
sudo mv skmp /usr/local/bin/
```

See the [Download page]({{ site.baseurl }}/download/) for all platforms and methods.

---

## Launch the TUI

```bash
skmp
```

You'll see four tabs:

| Tab | Key | Contents |
|-----|-----|----------|
| Skills | `1` | All available skills in the registry |
| Bundles | `2` | Curated skill collections |
| My Skills | `3` | Your installed skills |
| My Bundles | `4` | Bundles you've partially or fully installed |

---

## Install your first skill

Navigate with `j`/`k` (or mouse scroll), press `i` to install.

Or from the command line:

```bash
skmp add caveman
```

Multiple at once:

```bash
skmp add caveman diagnose tdd
```

Install a whole bundle:

```bash
skmp add --bundle nitesh000/core
```

---

## Remove a skill

Select and press `x`, or:

```bash
skmp remove caveman
```

---

## Per-skill harness control

By default, skills are linked into every detected agent. To control access per agent:

1. Select the skill
2. Press `tab` to focus the detail pane
3. Use `j`/`k` to navigate harnesses
4. Press `space` or `enter` to toggle

You can also click harness rows directly with your mouse.

---

## Sync after installing a new agent

```bash
skmp sync
```

Re-links all existing skills into the newly detected harness.

---

## Refresh the registry

The registry is cached locally for 24 hours. Force a refresh:

```bash
skmp update
```

---

## CLI reference

| Command | Description |
|---------|-------------|
| `skmp` | Open the TUI |
| `skmp add <skill> [skill...]` | Install skills by name |
| `skmp add --bundle <name>` | Install all skills in a bundle |
| `skmp remove <skill> [skill...]` | Remove skills |
| `skmp list` | List installed skills |
| `skmp sync` | Re-link skills into newly installed harnesses |
| `skmp update` | Force-refresh the local registry cache |
| `skmp version` | Print the installed version |

---

## Supported harnesses

| Harness | Skills path | Detection |
|---------|-------------|-----------|
| pi | `~/.agents/skills/` | `pi` in PATH |
| Claude Code | `~/.claude/skills/` | `claude` in PATH |
| Antigravity IDE | `~/.gemini/antigravity-ide/skills/` | `agy-ide` binary exists |
| agy | `~/.gemini/config/skills/` | `agy` in PATH |
| Codex | `~/.codex/skills/` | `codex` in PATH |
| Cursor | `~/.cursor/skills-cursor/` | `cursor` in PATH |
| OpenCode | registered in `opencode.jsonc` | `opencode` in PATH |

---

## Submit a skill

1. Create a repo with at least a `SKILL.md`
2. Open a PR adding your skill to [`registry/index.json`](https://github.com/Nitesh000/skmp/blob/master/registry/index.json)

See [CONTRIBUTING.md](https://github.com/Nitesh000/skmp/blob/master/CONTRIBUTING.md) for full details.

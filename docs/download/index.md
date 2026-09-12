---
layout: default
title: Download
nav_order: 3
description: Download skmp for macOS, Linux, or Windows.
---

# Download
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## npm (all platforms)

The easiest way to install on any platform with Node.js:

```bash
npm install -g skmp-cli
```

---

## Homebrew (macOS / Linux)

```bash
brew tap Nitesh000/tap
brew install skmp
```

---

## Binary downloads

Download pre-built binaries from [GitHub Releases](https://github.com/Nitesh000/skmp/releases/latest).

### macOS

```bash
# Apple Silicon (M1/M2/M3/M4)
curl -L https://github.com/Nitesh000/skmp/releases/latest/download/skmp_Darwin_arm64.tar.gz | tar xz
sudo mv skmp /usr/local/bin/

# Intel
curl -L https://github.com/Nitesh000/skmp/releases/latest/download/skmp_Darwin_x86_64.tar.gz | tar xz
sudo mv skmp /usr/local/bin/
```

### Linux

```bash
# x86_64
curl -L https://github.com/Nitesh000/skmp/releases/latest/download/skmp_Linux_x86_64.tar.gz | tar xz
sudo mv skmp /usr/local/bin/

# ARM64
curl -L https://github.com/Nitesh000/skmp/releases/latest/download/skmp_Linux_arm64.tar.gz | tar xz
sudo mv skmp /usr/local/bin/
```

### Windows

Download the `.zip` from [GitHub Releases](https://github.com/Nitesh000/skmp/releases/latest), extract, and add to your PATH.

---

## Verify installation

```bash
skmp version
```

---

## All releases

See the full [release history on GitHub](https://github.com/Nitesh000/skmp/releases) for older versions and release notes.

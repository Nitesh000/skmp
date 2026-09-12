package harness

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var SkillsDir = func() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".agents", "skills")
}

var StoreDir = func() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".skmp", "skills")
}

func InstalledSkills() ([]string, error) {
	dir := StoreDir()
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(dir, e.Name(), "SKILL.md")); err == nil {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func IsInstalled(name string) bool {
	_, err := os.Stat(filepath.Join(StoreDir(), name, "SKILL.md"))
	return err == nil
}

func Install(name, sourceURL string) error {
	storeDir := filepath.Join(StoreDir(), name)
	if err := os.MkdirAll(storeDir, 0755); err != nil {
		return fmt.Errorf("create store dir: %w", err)
	}

	files := []string{"SKILL.md", "REFERENCE.md", "EXAMPLES.md"}
	for _, f := range files {
		downloadFile(sourceURL+f, filepath.Join(storeDir, f))
	}

	if _, err := os.Stat(filepath.Join(storeDir, "SKILL.md")); err != nil {
		os.RemoveAll(storeDir)
		return fmt.Errorf("SKILL.md not found at %s", sourceURL)
	}

	if err := registerOpenCodeIfPresent(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: opencode config: %s\n", err)
	}

	if runtime.GOOS == "windows" {
		return copyToHarnessDirs(name, storeDir)
	}
	return symlinkToHarnessDirs(name, storeDir)
}

func Remove(name string) error {
	store := StoreDir()
	for _, dir := range allLinkTargets() {
		p := filepath.Join(dir, name)
		if !pointsInto(p, store) {
			continue
		}
		os.RemoveAll(p)
	}
	return os.RemoveAll(filepath.Join(store, name))
}

// linkTargets returns deduplicated dirs for currently-installed harnesses,
// excluding ConfigBased ones (they read the store directly).
// Defined as a var so tests can override it without needing real harnesses.
var linkTargets = func() []string {
	return collectTargets(InstalledHarnesses())
}

// allLinkTargets returns dirs for every known harness, used by Remove and
// Uninstall so they clean up even if a harness was uninstalled after the
// skill was added.
func allLinkTargets() []string {
	return collectTargets(Detect())
}

func collectTargets(harnesses []Harness) []string {
	seen := map[string]bool{}
	var dirs []string
	for _, h := range harnesses {
		if h.ConfigBased || seen[h.SkillsDir] {
			continue
		}
		seen[h.SkillsDir] = true
		dirs = append(dirs, h.SkillsDir)
	}
	return dirs
}

func symlinkToHarnessDirs(name, storeDir string) error {
	for _, dir := range linkTargets() {
		os.MkdirAll(dir, 0755)
		link := filepath.Join(dir, name)
		// Remove a stale symlink (broken link) before attempting to create.
		if isSymlink(link) {
			if _, err := os.Stat(link); os.IsNotExist(err) {
				os.Remove(link)
			}
		}
		if err := os.Symlink(storeDir, link); err != nil {
			if !os.IsExist(err) {
				return fmt.Errorf("symlink %s: %w", dir, err)
			}
			// Existing path: only replace if it's a symlink we own.
			if !isSymlink(link) {
				continue
			}
			os.Remove(link)
			if err := os.Symlink(storeDir, link); err != nil {
				return fmt.Errorf("symlink %s: %w", dir, err)
			}
		}
	}
	return nil
}

func copyToHarnessDirs(name, storeDir string) error {
	for _, dir := range linkTargets() {
		os.MkdirAll(dir, 0755)
		dest := filepath.Join(dir, name)
		if err := copyDir(storeDir, dest); err != nil {
			return fmt.Errorf("copy to %s: %w", dir, err)
		}
	}
	return nil
}

func copyDir(src, dst string) error {
	os.MkdirAll(dst, 0755)
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		sp := filepath.Join(src, e.Name())
		dp := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDir(sp, dp); err != nil {
				return err
			}
			continue
		}
		data, err := os.ReadFile(sp)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dp, data, 0644); err != nil {
			return err
		}
	}
	return nil
}

func downloadFile(url, dest string) error {
	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", res.StatusCode)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	return err
}

// pointsInto returns true when p is a symlink whose target lives under root,
// or on Windows when p is a directory copy managed by skmp.
func pointsInto(p, root string) bool {
	if runtime.GOOS == "windows" {
		// On Windows we copy, so just check if the directory exists and has a
		// SKILL.md — we can't distinguish our copies from user dirs, but the
		// matching store entry is also being deleted so this is safe.
		_, err := os.Stat(filepath.Join(p, "SKILL.md"))
		return err == nil
	}
	target, err := os.Readlink(p)
	if err != nil {
		return false
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return false
	}
	return strings.HasPrefix(abs, root+string(filepath.Separator))
}

func isSymlink(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink != 0
}

// Uninstall removes all links skmp created and deletes ~/.skmp.
func Uninstall() error {
	home, _ := os.UserHomeDir()
	store := StoreDir()

	entries, err := os.ReadDir(store)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	for _, dir := range allLinkTargets() {
		for _, e := range entries {
			p := filepath.Join(dir, e.Name())
			if !pointsInto(p, store) {
				continue
			}
			os.RemoveAll(p)
		}
	}

	unregisterOpenCode()

	return os.RemoveAll(filepath.Join(home, ".skmp"))
}

// --- opencode config management ---

var opencodeCfgPath = func() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "opencode", "opencode.jsonc")
}

func registerOpenCodeIfPresent() error {
	for _, h := range InstalledHarnesses() {
		if h.Name == "opencode" {
			return registerOpenCode()
		}
	}
	return nil
}

func registerOpenCode() error {
	cfgPath := opencodeCfgPath()
	data, err := os.ReadFile(cfgPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	store := StoreDir()
	var cfg map[string]any
	if err := json.Unmarshal(stripJSONC(data), &cfg); err != nil {
		return fmt.Errorf("parse opencode config: %w", err)
	}

	skills, _ := cfg["skills"].(map[string]any)
	if skills == nil {
		skills = map[string]any{}
	}
	paths, _ := skills["paths"].([]any)
	for _, p := range paths {
		if s, ok := p.(string); ok && s == store {
			return nil
		}
	}

	skills["paths"] = append(paths, store)
	cfg["skills"] = skills

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfgPath, out, 0644)
}

func unregisterOpenCode() {
	cfgPath := opencodeCfgPath()
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return
	}

	store := StoreDir()
	var cfg map[string]any
	if err := json.Unmarshal(stripJSONC(data), &cfg); err != nil {
		return
	}

	skills, _ := cfg["skills"].(map[string]any)
	if skills == nil {
		return
	}
	paths, _ := skills["paths"].([]any)
	var filtered []any
	for _, p := range paths {
		if s, ok := p.(string); ok && s == store {
			continue
		}
		filtered = append(filtered, p)
	}
	if len(filtered) == len(paths) {
		return
	}

	skills["paths"] = filtered
	cfg["skills"] = skills

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(cfgPath, out, 0644)
}

// stripJSONC removes // and /* */ comments from JSONC, respecting strings.
func stripJSONC(data []byte) []byte {
	var out []byte
	i := 0
	for i < len(data) {
		if data[i] == '"' {
			j := i + 1
			for j < len(data) {
				if data[j] == '\\' {
					j += 2
					continue
				}
				if data[j] == '"' {
					j++
					break
				}
				j++
			}
			out = append(out, data[i:j]...)
			i = j
		} else if i+1 < len(data) && data[i] == '/' && data[i+1] == '/' {
			for i < len(data) && data[i] != '\n' {
				i++
			}
		} else if i+1 < len(data) && data[i] == '/' && data[i+1] == '*' {
			i += 2
			for i+1 < len(data) && !(data[i] == '*' && data[i+1] == '/') {
				i++
			}
			if i+1 < len(data) {
				i += 2
			}
		} else {
			out = append(out, data[i])
			i++
		}
	}
	return out
}

func Sync() (int, error) {
	storeDir := StoreDir()
	entries, err := os.ReadDir(storeDir)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	if err := registerOpenCodeIfPresent(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: opencode config: %s\n", err)
	}

	uniqueDirs := linkTargets()
	linked := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		skillStore := filepath.Join(storeDir, name)

		for _, dir := range uniqueDirs {
			link := filepath.Join(dir, name)
			if _, err := os.Lstat(link); err == nil {
				continue
			}
			os.MkdirAll(dir, 0755)
			if runtime.GOOS == "windows" {
				copyDir(skillStore, link)
			} else {
				os.Symlink(skillStore, link)
			}
			linked++
		}
	}

	return linked, nil
}

// returns map of harness name: true, if the skill directory exist inside that harness's SkillDir
func SkillHarnessState(skillName string) map[string]bool {
	state := map[string]bool{}
	for _, h := range Detect() {
		_, err := os.Stat(filepath.Join(h.SkillsDir, skillName))
		state[h.Name] = err == nil
	}
	return state
}

func LinkSkillTo(skillName, harnessName string) error {
	storeEntry := filepath.Join(StoreDir(), skillName)
	for _, h := range Detect() {
		if h.Name != harnessName {
			continue
		}
		if h.ConfigBased {
			return nil
		}
		os.MkdirAll(h.SkillsDir, 0755)
		link := filepath.Join(h.SkillsDir, skillName)
		if runtime.GOOS == "windows" {
			return copyDir(storeEntry, skillName)
		}
		if isSymlink(link) {
			if _, err := os.Stat(link); os.IsNotExist(err) {
				os.Remove(link)
			}
		}
		if err := os.Symlink(storeEntry, link); err != nil {
			if !os.IsExist(err) {
				return fmt.Errorf("symlink %s: %w", h.SkillsDir, err)
			}
			if !isSymlink(link) {
				return nil
			}
			os.Remove(link)
			return os.Symlink(storeEntry, link)
		}
		return nil
	}
	return fmt.Errorf("harness %q not fouond", harnessName)
}

// remove skill from the harness's skilldir
func UnlinkSkillFrom(skillName, harnessName string) error {
	store := StoreDir()
	for _, h := range Detect() {
		if h.Name != harnessName {
			continue
		}
		if h.ConfigBased {
			return nil
		}
		p := filepath.Join(h.SkillsDir, skillName)
		if !pointsInto(p, store) {
			return nil
		}
		return os.RemoveAll(p)
	}
	return fmt.Errorf("harness %q not found", harnessName)
}

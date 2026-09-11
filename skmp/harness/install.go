package harness

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func SkillsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".agents", "skills")
}

func StoreDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".skmp", "skills")
}

// get list if installed skills
func InstalledSkills() ([]string, error) {
	dir := SkillsDir()

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
		skillMD := filepath.Join(dir, e.Name(), "SKILL.md")
		if _, err := os.Stat(skillMD); err == nil {
			names = append(names, e.Name())
		}
	}

	return names, nil
}

// check if a skill is installed or not
func IsInstalled(name string) bool {
	p := filepath.Join(SkillsDir(), name, "SKILL.md")
	_, err := os.Stat(p)
	return err == nil
}

// install skill
func Install(name, sourceURL string) error {
	// download the files into ~/.skmp/skills/<name>
	storeDir := filepath.Join(StoreDir(), name)
	if err := os.MkdirAll(storeDir, 0755); err != nil {
		return fmt.Errorf("create store dir: %w", err)
	}

	files := []string{"SKILL.md", "REFERENCE.md", "EXAMPLES.md"}
	downloaded := 0
	for _, f := range files {
		url := sourceURL + f
		dest := filepath.Join(storeDir, f)
		if err := downloadFile(url, dest); err == nil {
			downloaded++
		}
	}

	if downloaded == 0 {
		os.RemoveAll(storeDir)
		return fmt.Errorf("no files downloaded from %s", sourceURL)
	}

	for _, h := range InstalledHarnesses() {
		if h.Name == "opencode" {
			registerOpenCode()
		}
	}

	// link into harness dir
	if runtime.GOOS == "windows" {
		return copyToHarnessDirs(name, storeDir)
	}

	return symlinkToHarnessDirs(name, storeDir)
}

// remove skill
func Remove(name string) error {
	for _, dir := range linkTargets() {
		os.RemoveAll(filepath.Join(dir, name))
	}

	// remove from local store
	return os.RemoveAll(filepath.Join(StoreDir(), name))
}

// linkTargets returns every directory that needs a link, deduplicated.
// Config-based harnesses are excluded — they read from the store directly.
func linkTargets() []string {
	seen := map[string]bool{}
	var dirs []string
	for _, h := range InstalledHarnesses() {
		if h.ConfigBased || seen[h.SkillsDir] {
			continue
		}
		seen[h.SkillsDir] = true
		dirs = append(dirs, h.SkillsDir)
	}
	return dirs
}

// generating symlinks
func symlinkToHarnessDirs(name, storeDir string) error {
	for _, dir := range linkTargets() {
		os.MkdirAll(dir, 0755)
		link := filepath.Join(dir, name)
		if err := os.Symlink(storeDir, link); err != nil && !os.IsExist(err) {
			return fmt.Errorf("symlink %s: %w", dir, err)
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
	entries, _ := os.ReadDir(src)
	for _, e := range entries {
		data, err := os.ReadFile(filepath.Join(src, e.Name()))
		if err != nil {
			continue
		}
		os.WriteFile(filepath.Join(dst, e.Name()), data, 0644)
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

// Uninstall removes all symlinks skmp created and deletes ~/.skmp entirely
func Uninstall() error {
	home, _ := os.UserHomeDir()
	storeDir := StoreDir()

	// remove all symlinks pointing into ~/.skmp/skills from every harness dir
	entries, err := os.ReadDir(storeDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	seen := map[string]bool{}
	for _, h := range Detect() {
		if h.ConfigBased || seen[h.SkillsDir] {
			continue
		}
		seen[h.SkillsDir] = true
		for _, e := range entries {
			os.Remove(filepath.Join(h.SkillsDir, e.Name()))
		}
	}

	// delete ~/.skmp entirely
	return os.RemoveAll(filepath.Join(home, ".skmp"))
}

// registerOpenCode adds the skmp store to opencode's skills.paths config.
// OpenCode reads skills from paths listed in its config rather than a fixed dir.
func registerOpenCode() error {
	home, _ := os.UserHomeDir()
	cfgPath := filepath.Join(home, ".config", "opencode", "opencode.jsonc")

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return err
	}

	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	store := StoreDir()
	skills, _ := cfg["skills"].(map[string]any)
	if skills == nil {
		skills = map[string]any{}
	}

	paths, _ := skills["paths"].([]any)
	for _, p := range paths {
		if s, ok := p.(string); ok && s == store {
			return nil // already registered
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

func Sync() (int, error) {
	storeDir := StoreDir()
	entries, err := os.ReadDir(storeDir)

	if os.IsNotExist(err) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	for _, h := range InstalledHarnesses() {
		if h.Name == "opencode" {
			registerOpenCode()
		}
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

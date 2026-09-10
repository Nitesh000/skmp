package harness

import (
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

	// link into harness dir
	if runtime.GOOS == "windows" {
		return copyToHarnessDirs(name, storeDir)
	}

	return symlinkToHarnessDirs(name, storeDir)
}

// remove skill
func Remove(name string) error {
	// remove symlink/copies
	os.RemoveAll(filepath.Join(SkillsDir(), name))

	home, _ := os.UserHomeDir()
	codexSkill := filepath.Join(home, ".codex", "skills", name)
	os.RemoveAll(codexSkill)

	// remove from local store
	return os.RemoveAll(filepath.Join(StoreDir(), name))
}

// generating symlinks
func symlinkToHarnessDirs(name, storeDir string) error {
	// ~/.agents directory
	agentsLink := filepath.Join(SkillsDir(), name)
	os.MkdirAll(filepath.Dir(agentsLink), 0755)
	if err := os.Symlink(storeDir, agentsLink); err != nil && !os.IsExist(err) {
		return fmt.Errorf("symlink agents: %w", err)
	}

	// ~/.codex directory (only if it's installed)
	home, _ := os.UserHomeDir()
	codexDir := filepath.Join(home, ".codex")
	if _, err := os.Stat(codexDir); err == nil {
		codexLink := filepath.Join(codexDir, "skills", name)
		os.MkdirAll(filepath.Dir(codexLink), 0755)
		os.Symlink(storeDir, codexLink)
	}

	return nil
}

func copyToHarnessDirs(name, storeDir string) error {
	targets := []string{filepath.Join(SkillsDir(), name)}

	home, _ := os.UserHomeDir()
	codexDir := filepath.Join(home, ".codex")
	if _, err := os.Stat(codexDir); err == nil {
		targets = append(targets, filepath.Join(codexDir, "skills", name))
	}

	for _, t := range targets {
		if err := copyDir(storeDir, t); err != nil {
			return err
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

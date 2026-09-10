package harness

import (
	"os"
	"path/filepath"
)

func SkillsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".agents", "skills")
}

func StoreDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".skmp", "skills")
}

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

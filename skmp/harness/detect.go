package harness

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Harness struct {
	Name        string
	SkillsDir   string
	Installed   bool
	ConfigBased bool
}

func Detect() []Harness {
	home, _ := os.UserHomeDir()
	agentsSkills := SkillsDir()

	return []Harness{
		{
			Name:      "pi",
			SkillsDir: agentsSkills,
			Installed: commandExists("pi"),
		},
		{
			Name:      "claude-code",
			SkillsDir: filepath.Join(home, ".claude", "skills"),
			Installed: commandExists("claude"),
		},
		{
			Name:      "antigravity-ide",
			SkillsDir: filepath.Join(home, ".gemini", "antigravity-ide", "skills"),
			Installed: fileExists(filepath.Join(home, ".antigravity-ide", "antigravity-ide", "bin", "agy-ide")),
		},
		{
			Name:      "agy",
			SkillsDir: filepath.Join(home, ".gemini", "config", "skills"),
			Installed: commandExists("agy"),
		},
		{
			Name:        "opencode",
			SkillsDir:   StoreDir(),
			Installed:   commandExists("opencode"),
			ConfigBased: true,
		},
		{
			Name:      "codex",
			SkillsDir: filepath.Join(home, ".codex", "skills"),
			Installed: commandExists("codex"),
		},
		{
			Name:      "cursor",
			SkillsDir: filepath.Join(home, ".cursor", "skills-cursor"),
			Installed: commandExists("cursor"),
		},
	}
}

func InstalledHarnesses() []Harness {
	var out []Harness
	for _, h := range Detect() {
		if h.Installed {
			out = append(out, h)
		}
	}
	return out
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

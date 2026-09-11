package harness

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Harness struct {
	Name      string
	SkillsDir string
	Installed bool
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
			SkillsDir: agentsSkills,
			Installed: commandExists("claude"),
		},
		{
			Name:      "opencode",
			SkillsDir: agentsSkills,
			Installed: commandExists("opencode"),
		},
		{
			Name:      "antigravity",
			SkillsDir: agentsSkills,
			Installed: dirExist(filepath.Join(home, ".antigravity-ide")),
		},
		{
			Name:      "codex",
			SkillsDir: filepath.Join(home, ".codex", "skills"),
			Installed: commandExists("codex"),
		},
		{
			Name:      "cursor",
			SkillsDir: agentsSkills,
			Installed: dirExist(filepath.Join(home, ".cursor")),
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

func dirExist(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

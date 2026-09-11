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
	// ConfigBased harnesses read skills from a path registered in their own
	// config file instead of scanning a well-known directory.
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
			SkillsDir: agentsSkills,
			Installed: commandExists("claude"),
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

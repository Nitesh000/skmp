package harness

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func setupTestDirs(t *testing.T) (store, agents string) {
	t.Helper()
	tmp := t.TempDir()
	store = filepath.Join(tmp, "store")
	agents = filepath.Join(tmp, "agents")
	os.MkdirAll(store, 0755)
	os.MkdirAll(agents, 0755)
	return
}

func writeSkill(t *testing.T, dir, name string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	os.MkdirAll(p, 0755)
	os.WriteFile(filepath.Join(p, "SKILL.md"), []byte("# "+name), 0644)
	return p
}

// --- Bug 1: source of truth ---

func TestInstalledSkillsReadsStore(t *testing.T) {
	store, agents := setupTestDirs(t)

	// Override the dirs for the test.
	origStore := StoreDir
	origSkills := SkillsDir
	StoreDir = func() string { return store }
	SkillsDir = func() string { return agents }
	defer func() { StoreDir = origStore; SkillsDir = origSkills }()

	// Put a skill only in agents (user-owned, not skmp).
	writeSkill(t, agents, "user-owned")

	// Put a skill in the store (skmp-managed).
	writeSkill(t, store, "managed")

	skills, err := InstalledSkills()
	if err != nil {
		t.Fatal(err)
	}

	if len(skills) != 1 || skills[0] != "managed" {
		t.Errorf("expected [managed], got %v", skills)
	}

	if IsInstalled("user-owned") {
		t.Error("user-owned should not be reported as installed")
	}
	if !IsInstalled("managed") {
		t.Error("managed should be reported as installed")
	}
}

func TestRemoveOnlyUnlinksOwnedSymlinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink test")
	}

	store, agents := setupTestDirs(t)

	origStore := StoreDir
	origSkills := SkillsDir
	StoreDir = func() string { return store }
	SkillsDir = func() string { return agents }
	defer func() { StoreDir = origStore; SkillsDir = origSkills }()

	// A user-owned directory sitting in agents.
	userDir := writeSkill(t, agents, "user-owned")

	// A skmp-managed skill: store entry + symlink.
	storeEntry := writeSkill(t, store, "managed")
	os.Symlink(storeEntry, filepath.Join(agents, "managed"))

	if err := Remove("user-owned"); err != nil {
		t.Fatal(err)
	}
	// user-owned must survive.
	if _, err := os.Stat(userDir); err != nil {
		t.Error("Remove deleted a user-owned directory")
	}

	if err := Remove("managed"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(agents, "managed")); err == nil {
		t.Error("symlink should have been removed")
	}
	if _, err := os.Stat(storeEntry); err == nil {
		t.Error("store entry should have been removed")
	}
}

// --- Bug 2: registerOpenCode JSONC ---

func TestRegisterOpenCodeStripsComments(t *testing.T) {
	tmp := t.TempDir()
	cfgDir := filepath.Join(tmp, ".config", "opencode")
	os.MkdirAll(cfgDir, 0755)
	cfgPath := filepath.Join(cfgDir, "opencode.jsonc")

	origCfgPath := opencodeCfgPath
	opencodeCfgPath = func() string { return cfgPath }
	defer func() { opencodeCfgPath = origCfgPath }()

	jsonc := `{
  // this is a comment
  "$schema": "https://opencode.ai/config.json"
  /* block comment */
}`
	os.WriteFile(cfgPath, []byte(jsonc), 0644)

	store := "/tmp/fake-store"
	origStore := StoreDir
	StoreDir = func() string { return store }
	defer func() { StoreDir = origStore }()

	if err := registerOpenCode(); err != nil {
		t.Fatalf("registerOpenCode failed on JSONC: %v", err)
	}

	data, _ := os.ReadFile(cfgPath)
	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	skills := cfg["skills"].(map[string]any)
	paths := skills["paths"].([]any)
	if len(paths) != 1 || paths[0] != store {
		t.Errorf("expected [%s], got %v", store, paths)
	}
}

func TestRegisterOpenCodeMissingFile(t *testing.T) {
	origCfgPath := opencodeCfgPath
	opencodeCfgPath = func() string { return filepath.Join(t.TempDir(), "nope.jsonc") }
	defer func() { opencodeCfgPath = origCfgPath }()

	if err := registerOpenCode(); err != nil {
		t.Errorf("missing file should be a no-op, got: %v", err)
	}
}

func TestUnregisterOpenCode(t *testing.T) {
	tmp := t.TempDir()
	cfgDir := filepath.Join(tmp, ".config", "opencode")
	os.MkdirAll(cfgDir, 0755)
	cfgPath := filepath.Join(cfgDir, "opencode.jsonc")

	origCfgPath := opencodeCfgPath
	opencodeCfgPath = func() string { return cfgPath }
	defer func() { opencodeCfgPath = origCfgPath }()

	store := "/tmp/fake-store"
	origStore := StoreDir
	StoreDir = func() string { return store }
	defer func() { StoreDir = origStore }()

	cfg := map[string]any{
		"skills": map[string]any{
			"paths": []any{store, "/other"},
		},
	}
	data, _ := json.Marshal(cfg)
	os.WriteFile(cfgPath, data, 0644)

	unregisterOpenCode()

	after, _ := os.ReadFile(cfgPath)
	var out map[string]any
	json.Unmarshal(after, &out)
	paths := out["skills"].(map[string]any)["paths"].([]any)
	if len(paths) != 1 || paths[0] != "/other" {
		t.Errorf("expected [/other], got %v", paths)
	}
}

// --- Bug 3: stale symlinks ---

func TestSymlinkReplacesStaleLink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink test")
	}

	store, agents := setupTestDirs(t)

	origStore := StoreDir
	origSkills := SkillsDir
	StoreDir = func() string { return store }
	SkillsDir = func() string { return agents }
	defer func() { StoreDir = origStore; SkillsDir = origSkills }()

	storeEntry := writeSkill(t, store, "skill-a")

	// Create a stale symlink pointing nowhere.
	staleTarget := filepath.Join(t.TempDir(), "old-location")
	os.Symlink(staleTarget, filepath.Join(agents, "skill-a"))

	if err := symlinkToHarnessDirs("skill-a", storeEntry); err != nil {
		t.Fatal(err)
	}

	target, err := os.Readlink(filepath.Join(agents, "skill-a"))
	if err != nil {
		t.Fatal(err)
	}
	if target != storeEntry {
		t.Errorf("expected link to %s, got %s", storeEntry, target)
	}
}

func TestSymlinkPreservesRealDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink test")
	}

	store, agents := setupTestDirs(t)

	origStore := StoreDir
	origSkills := SkillsDir
	StoreDir = func() string { return store }
	SkillsDir = func() string { return agents }
	defer func() { StoreDir = origStore; SkillsDir = origSkills }()

	storeEntry := writeSkill(t, store, "skill-a")
	writeSkill(t, agents, "skill-a") // real user dir

	if err := symlinkToHarnessDirs("skill-a", storeEntry); err != nil {
		t.Fatal(err)
	}

	// The real dir must survive — it's not a symlink so we skip it.
	fi, err := os.Lstat(filepath.Join(agents, "skill-a"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		t.Error("real user directory was replaced with a symlink")
	}
}

// --- Bug 6: SKILL.md required ---

func TestInstallFailsWithoutSkillMD(t *testing.T) {
	store := filepath.Join(t.TempDir(), "store")
	origStore := StoreDir
	StoreDir = func() string { return store }
	defer func() { StoreDir = origStore }()

	// Use a source URL that 404s everything.
	err := Install("fake-skill", "http://127.0.0.1:1/nonexistent/")
	if err == nil {
		t.Fatal("expected error when nothing downloads")
	}
	if !strings.Contains(err.Error(), "SKILL.md") {
		t.Errorf("error should mention SKILL.md, got: %s", err)
	}

	// Store dir must be cleaned up.
	if _, err := os.Stat(filepath.Join(store, "fake-skill")); err == nil {
		t.Error("store dir should have been cleaned up")
	}
}

// --- Bug 7: recursive copyDir ---

func TestCopyDirRecursive(t *testing.T) {
	src := filepath.Join(t.TempDir(), "src")
	dst := filepath.Join(t.TempDir(), "dst")

	os.MkdirAll(filepath.Join(src, "sub"), 0755)
	os.WriteFile(filepath.Join(src, "top.md"), []byte("top"), 0644)
	os.WriteFile(filepath.Join(src, "sub", "nested.md"), []byte("nested"), 0644)

	if err := copyDir(src, dst); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dst, "sub", "nested.md"))
	if err != nil {
		t.Fatal("nested file not copied:", err)
	}
	if string(data) != "nested" {
		t.Errorf("content mismatch: %q", data)
	}
}

// --- Bug 5: Uninstall on Windows (os.RemoveAll) ---

func TestPointsIntoRejectsNonSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix-only")
	}

	store, agents := setupTestDirs(t)
	writeSkill(t, agents, "user-owned")

	if pointsInto(filepath.Join(agents, "user-owned"), store) {
		t.Error("a real directory should not be considered managed by skmp")
	}
}

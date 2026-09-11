package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	indexURL = "https://cdn.jsdelivr.net/gh/nitesh000/skmp@master/registry/index.json"
	cacheTTL = 24 * time.Hour
)

func cachePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".skmp", "cache", "index.json")
}

func cacheValid() bool {
	info, err := os.Stat(cachePath())
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) < cacheTTL
}

func fetchRemote() (*Index, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(indexURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", res.StatusCode)
	}

	var idx Index
	if err := json.NewDecoder(res.Body).Decode(&idx); err != nil {
		return nil, err
	}
	return &idx, nil
}

func readCache() (*Index, error) {
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return nil, err
	}
	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return &idx, nil
}

func saveCache(idx *Index) error {
	path := cachePath()
	os.MkdirAll(filepath.Dir(path), 0755)
	data, _ := json.Marshal(idx)
	return os.WriteFile(path, data, 0644)
}

func Load() (*Index, error) {
	if cacheValid() {
		return readCache()
	}

	return UpdateCache()
}

// UpdateCache ignores TTL and forces a network fetch
func UpdateCache() (*Index, error) {
	idx, err := fetchRemote()
	if err != nil {
		// network failed — fall back to stale cache if it exists
		if cached, cacheErr := readCache(); cacheErr == nil {
			return cached, nil
		}
		return nil, fmt.Errorf("fetch failed and no cache: %w", err)
	}

	saveCache(idx)
	return idx, nil
}

func LoadBundle(name string, idx *Index) (*BundleFile, error) {
	for _, b := range idx.Bundles {
		if b.Name == name {
			return &BundleFile{
				Bundle:      b.Name,
				Description: b.Description,
				Author:      b.Author,
				Repo:        b.Repo,
				Branch:      b.Branch,
				SkillsPath:  b.SkillsPath,
				Skills:      b.Skills,
			}, nil
		}
	}
	return nil, fmt.Errorf("bundle %q not found in registry", name)
}

// BundleSkillSource builds the download URL for a skill inside a bundle.
// Used when a bundle skill is not individually listed in index.json.
func BundleSkillSource(bf *BundleFile, skillName string) string {
	rawBase := strings.Replace(bf.Repo, "https://github.com/", "https://raw.githubusercontent.com/", 1)
	rawBase = strings.TrimRight(rawBase, "/")
	skillsPath := strings.Trim(bf.SkillsPath, "/")
	return rawBase + "/" + bf.Branch + "/" + skillsPath + "/" + skillName + "/"
}

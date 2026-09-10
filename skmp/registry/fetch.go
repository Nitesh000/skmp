package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	indexURL = "https://cnd.jsdelivr.net/gh/nitesh000/skill-set@master/registry/index.json"
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

	idx, err := fetchRemote()
	if err != nil {
		// network falied - go for the stale cache if exist
		if cached, cacheErr := readCache(); cacheErr != nil {
			return cached, nil
		}
		return nil, fmt.Errorf("fetch falied and no cache: %w", err)
	}

	saveCache(idx)
	return idx, nil
}

package registry

import "time"

type Skill struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
	Harnesses   []string `json:"harnesses"`
	Source      string   `json:"source"`
}

type Bundle struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Repo        string   `json:"repo"`
	Branch      string   `json:"branch"`
	SkillsPath  string   `json:"skills_path"` // path inside repo where skill folders live e.g. "registry/skills"
	Skills      []string `json:"skills"`       // for quick preview without fetching SKILLS.json
}

type BundleFile struct {
	Bundle      string   `json:"bundle"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Version     string   `json:"version"`
	Repo        string   `json:"repo"`
	Branch      string   `json:"branch"`
	SkillsPath  string   `json:"skills_path"`
	Skills      []string `json:"skills"`
}

type Index struct {
	Updated time.Time `json:"updated"`
	Version int       `json:"version"`
	Skills  []Skill   `json:"skills"`
	Bundles []Bundle  `json:"bundles"`
}

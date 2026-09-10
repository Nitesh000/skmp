package registry

import "time"

type Skill struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Tags        []string `json:"Tags"`
	Harnesses   []string `json:"harness"`
	Source      string   `json:"source"`
}

type Index struct {
	Updated time.Time `json:"updated"`
	Version int       `json:"version"`
	Skills  []Skill   `json:"skills"`
}

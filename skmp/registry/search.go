package registry

import (
	"strings"

	bleve "github.com/blevesearch/bleve/v2"
)

var (
	index       bleve.Index
	allSkills   []Skill
	skillByName map[string]Skill
)

func BuildIndex(skills []Skill) error {
	allSkills = skills
	skillByName = map[string]Skill{}

	mapping := bleve.NewIndexMapping()
	idx, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return err
	}

	batch := idx.NewBatch()
	for _, s := range skills {
		doc := map[string]string{
			"name":        s.Name,
			"description": s.Description,
			"tags":        strings.Join(s.Tags, " "),
			"author":      s.Author,
		}
		batch.Index(s.Name, doc)
		skillByName[s.Name] = s
	}

	if err := idx.Batch(batch); err != nil {
		return err
	}

	index = idx
	return nil
}

func Search(query string) ([]Skill, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return allSkills, nil
	}

	nameQ := bleve.NewFuzzyQuery(query)
	nameQ.Fuzziness = 1

	tagQ := bleve.NewMatchQuery(query)
	tagQ.SetField("tags")

	combined := bleve.NewDisjunctionQuery(nameQ, tagQ)
	req := bleve.NewSearchRequestOptions(combined, 50, 0, false)

	res, err := index.Search(req)
	if err != nil {
		return []Skill{}, nil
	}

	seen := map[string]struct{}{}
	var out []Skill

	// bleve hits in relevance order
	for _, h := range res.Hits {
		if s, ok := skillByName[h.ID]; ok {
			out = append(out, s)
			seen[h.ID] = struct{}{}
		}
	}

	return out, nil
}

func fallback(query string) []Skill {
	q := strings.ToLower(query)
	var out []Skill
	for _, s := range allSkills {
		if strings.Contains(strings.ToLower(s.Name), q) ||
			strings.Contains(strings.ToLower(s.Description), q) ||
			strings.Contains(strings.ToLower(strings.Join(s.Tags, " ")), q) {
			out = append(out, s)
		}
	}
	return out
}

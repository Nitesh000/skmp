package registry

import (
	"strings"

	bleve "github.com/blevesearch/bleve/v2"
)

var (
	index        bleve.Index
	bundleIndex  bleve.Index
	allSkills    []Skill
	allBundles   []Bundle
	skillByName  map[string]Skill
	bundleByName map[string]Bundle
)

func BuildIndex(skills []Skill, bundles []Bundle) error {
	allSkills = skills
	allBundles = bundles
	skillByName = map[string]Skill{}
	bundleByName = map[string]Bundle{}

	// --- Skills Index ---
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

	// --- Bundles Index ---
	bIdx, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return err
	}

	bBatch := bIdx.NewBatch()
	for _, b := range bundles {
		doc := map[string]string{
			"name":        b.Name,
			"description": b.Description,
			"author":      b.Author,
			"skills":      strings.Join(b.Skills, " "),
		}
		bBatch.Index(b.Name, doc)
		bundleByName[b.Name] = b
	}

	if err := bIdx.Batch(bBatch); err != nil {
		return err
	}
	bundleIndex = bIdx

	return nil
}

func SearchSkills(query string) ([]Skill, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return allSkills, nil
	}

	lowerQ := strings.ToLower(query)

	// 1. Substring match (catches "c" in "caveman", or "dd" in "tdd")
	wildcard := bleve.NewWildcardQuery("*" + lowerQ + "*")

	// 2. Typo match (catches "cavman" -> "caveman")
	fuzzy := bleve.NewFuzzyQuery(lowerQ)
	fuzzy.Fuzziness = 2

	// Combine them — if it matches EITHER wildcard OR fuzzy, include it
	combined := bleve.NewDisjunctionQuery(wildcard, fuzzy)
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

func SearchBundles(query string) ([]Bundle, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return allBundles, nil
	}

	lowerQ := strings.ToLower(query)

	wildcard := bleve.NewWildcardQuery("*" + lowerQ + "*")
	fuzzy := bleve.NewFuzzyQuery(lowerQ)
	fuzzy.Fuzziness = 2

	combined := bleve.NewDisjunctionQuery(wildcard, fuzzy)
	req := bleve.NewSearchRequestOptions(combined, 50, 0, false)

	res, err := bundleIndex.Search(req)
	if err != nil {
		return []Bundle{}, nil
	}

	var out []Bundle
	for _, h := range res.Hits {
		if b, ok := bundleByName[h.ID]; ok {
			out = append(out, b)
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

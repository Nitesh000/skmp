package registry

import (
	"strings"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/mapping"
)

const analyzerName = "lowercase_keyword"

var (
	index        bleve.Index
	bundleIndex  bleve.Index
	allSkills    []Skill
	allBundles   []Bundle
	skillByName  map[string]Skill
	bundleByName map[string]Bundle
)

// newIndexMapping returns an index mapping that uses a custom analyzer which
// lowercases but does NOT tokenize on punctuation. This means hyphens, dots,
// and underscores are preserved, so wildcard queries like *write-a* match
// write-a-skill without any string-scan fallback.
func newIndexMapping() (*mapping.IndexMappingImpl, error) {
	im := bleve.NewIndexMapping()

	err := im.AddCustomAnalyzer(analyzerName, map[string]interface{}{
		"type":      keyword.Name,
		"token_filters": []string{lowercase.Name},
	})
	if err != nil {
		return nil, err
	}

	fm := bleve.NewTextFieldMapping()
	fm.Analyzer = analyzerName

	dm := bleve.NewDocumentMapping()
	dm.AddFieldMappingsAt("name", fm)
	dm.AddFieldMappingsAt("description", fm)
	dm.AddFieldMappingsAt("tags", fm)
	dm.AddFieldMappingsAt("author", fm)
	dm.AddFieldMappingsAt("skills", fm)

	im.DefaultMapping = dm
	return im, nil
}
func BuildIndex(skills []Skill, bundles []Bundle) error {
	allSkills = skills
	allBundles = bundles
	skillByName = map[string]Skill{}
	bundleByName = map[string]Bundle{}

	im, err := newIndexMapping()
	if err != nil {
		return err
	}

	// --- Skills Index ---
	idx, err := bleve.NewMemOnly(im)
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
	bIdx, err := bleve.NewMemOnly(im)
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


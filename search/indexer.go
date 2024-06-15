package search

import "fiber-search-engine/db"

// Index is an in-memory inverted index. It maps tokens to url IDs.
type Index map[string][]string

// Add is a method of the Index struct that adds a slice of CrawledUrl documents to the index.
// Adds documents to the Index.
// It loops over the documents and for each one, it analyzes the URL, page title, page description, and headings.
// For each token produced by the analysis, it checks if the document ID is already in the index for that token.
// If the document ID is not already in the index for that token, it adds the ID to the index.
// If the document ID is already in the index for that token, it does not add the ID again.
//
// Parameters:
// docs []db.CrawledUrl: A slice of CrawledUrl documents to add to the index.
//
// This method does not return any values.
func (idx Index) Add(docs []db.CrawledUrl) {
	for _, doc := range docs {
		for _, token := range analyze(doc.Url + " " + doc.PageTitle + " " + doc.PageDescription + " " + doc.Headings) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.ID {
				// Don't add same ID twice.
				continue
			}
			idx[token] = append(ids, doc.ID)
		}
	}
}

package search

import "fiber-search-engine/db"

// in memory representation of our search index. Inverted index where the key is a word and the value is a list of URLs.
// Inverted index is a data structure used to create full text search. It maps content, such as words or numbers, to its locations in a database.
type Index map[string][]string

// Add adds the given documents to the index.
// It iterates over the documents and for each one, it analyzes the URL, page title, page description, and heading to get tokens.
// For each token, it adds the document's ID to the index under the token key.
// If the last ID under the token key is the same as the current document's ID, it skips the addition.
//
// Parameters:
//
//	docs: A slice of CrawledUrl documents to add to the index.
func (idx Index) Add(docs []db.CrawledUrl) {
	for _, doc := range docs {
		for _, token := range analyze(doc.Url + " " + doc.PageTitle + " " + doc.PageDescription + " " + doc.Heading) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.ID { // if the last id is the same as the current id, skip
				continue
			}

			idx[token] = append(ids, doc.ID)
		}
	}
}

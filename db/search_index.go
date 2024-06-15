package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type SearchIndex struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Value     string
	Urls      []CrawledUrl   `gorm:"many2many:token_urls;"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName is a function that returns the name of the database table associated with the SearchIndex struct.
// It returns the string "search_index".
//
// This function does not take any parameters.
//
// Returns:
// string: The name of the database table associated with the SearchIndex struct.
func (s *SearchIndex) TableName() string {
	return "search_index"
}

// Save is a method on the SearchIndex struct that saves the search index to the database.
// It takes a map of strings to slices of strings, where the key is the search term and the value is a slice of URLs.
// It also takes a slice of CrawledUrl objects that represent the URLs that have been crawled.
// It iterates over the search index map and creates a new SearchIndex object for each search term.
// It then associates the URLs with the search term and saves the search index to the database.
//
// Parameters:
// index map[string][]string: A map of strings to slices of strings, where the key is the search term and the value is a slice of URLs.
// crawledUrls []CrawledUrl: A slice of CrawledUrl objects that represent the URLs that have been crawled.
//
// Returns:
// error: An error object that describes an error that occurred during the method's execution.
func (s *SearchIndex) Save(index map[string][]string, crawledUrls []CrawledUrl) error {
	for value, ids := range index {
		newIndex := &SearchIndex{
			Value: value,
		}
		if err := DBConn.Where(SearchIndex{Value: value}).FirstOrCreate(newIndex).Error; err != nil {
			return err
		}

		var urlsToAppend []CrawledUrl
		for _, id := range ids {
			for _, url := range crawledUrls {
				if url.ID == id {
					urlsToAppend = append(urlsToAppend, url)
					break
				}
			}
		}

		if err := DBConn.Model(&newIndex).Association("Urls").Append(&urlsToAppend); err != nil {
			return err
		}
	}
	return nil
}

// FullTextSearch is a method on the SearchIndex struct that performs a full-text search on the search index.
// It takes a string value as input and splits it into individual terms.
// It then retrieves all search indexes that contain any of the terms and retrieves the associated URLs.
// It returns a slice of CrawledUrl objects that match the search terms.
//
// Parameters:
// value string: The search query string.
//
// Returns:
// []CrawledUrl: A slice of CrawledUrl objects that match the search query.
// error: An error object that describes an error that occurred during the method's execution.
func (s *SearchIndex) FullTextSearch(value string) ([]CrawledUrl, error) {
	terms := strings.Fields(value)
	var urls []CrawledUrl

	for _, term := range terms {
		var searchIndexes []SearchIndex
		if err := DBConn.Preload("Urls").Where("value LIKE ?", "%"+term+"%").Find(&searchIndexes).Error; err != nil {
			// Preload is used to eager load the associated URLs with the search index
			return nil, err
		}

		for _, searchIndex := range searchIndexes {
			urls = append(urls, searchIndex.Urls...)
		}
	}
	return urls, nil
}

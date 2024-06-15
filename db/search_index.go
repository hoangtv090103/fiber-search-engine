package db

import (
	"time"

	"gorm.io/gorm"
)

type SearchIndex struct {
	ID       string `gorm:"type:uuid;default:uuid_generate_v4();"`
	Value    string
	Urls     []CrawledUrl   `gorm:"many2many:token_urls;"`
	CreateAt *time.Time     `gorm:"autoCreateTime"`
	UpdateAt time.Time      `gorm:"autoUpdateTime"`
	DeleteAt gorm.DeletedAt `gorm:"index"`
}

func (s *SearchIndex) TableName() string {
	return "search_index"
}

// Save persists the SearchIndex instance to the database.
// It iterates over the given index map and for each value, it creates a new SearchIndex instance and saves it to the database.
// It also associates the SearchIndex instance with the corresponding CrawledUrl instances.
// If there's an error during the save operation or the association, it returns the error.
//
// Parameters:
//
//	index: A map where the keys are string values and the values are slices of string IDs.
//	crawledUrls: A slice of CrawledUrl instances to associate with the SearchIndex instances.
//
// Returns:
//
//	error: An error that will be non-nil if there was an issue saving the instance or associating the URLs.
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

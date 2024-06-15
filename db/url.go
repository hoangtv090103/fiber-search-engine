package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CrawledUrl struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Url             string         `json:"url" gorm:"unique; not null"`
	Success         bool           `json:"success" gorm:"not null"`
	CrawlDuration   time.Duration  `json:"crawlDuration" gorm:"not null"`
	ResponseCode    int            `json:"responseCode" gorm:"not null"`
	PageTitle       string         `json:"pageTitle"`
	PageDescription string         `json:"pageDescription"`
	Heading         string         `json:"heading"`
	LastTested      *time.Time     `json:"lastTested"`
	Indexed         bool           `json:"indexed"`
	CreateAt        *time.Time     `gorm:"autoCreateTime"`
	UpdateAt        time.Time      `gorm:"autoUpdateTime"`
	DeleteAt        gorm.DeletedAt `gorm:"index"`
}

// UpdatedUrl updates the fields of the CrawledUrl instance in the database.
// It uses the global DBConn to save the instance.
// It selects the url, success, crawl_duration, response_code, page_title, page_description, headings, last_tested, and updated_at fields for updating.
// It omits the create_at field from updating.
// If there's an error during the save operation, it prints the error and returns it.
//
// Parameters:
//
//	input: The CrawledUrl instance to update.
//
// Returns:
//
//	error: An error that will be non-nil if there was an issue saving the instance.
func (crawled *CrawledUrl) UpdatedUrl(input CrawledUrl) error {
	tx := DBConn.Select("url", "success", "crawl_duration", "response_code", "page_title", "page_description", "headings", "last_tested", "updated_at").Omit("create_at").Save(&input)

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return tx.Error
	}
	return nil
}

// GetNextCrawlUrls retrieves the next set of URLs to crawl.
// It queries the database for CrawledUrl instances where the last_tested field is null, up to the given limit.
// It returns a slice of CrawledUrl instances and an error.
// If there's an error during the query, it prints the error and returns an empty slice and the error.
//
// Parameters:
//
//	limit: The maximum number of URLs to retrieve.
//
// Returns:
//
//	[]CrawledUrl: A slice of CrawledUrl instances to crawl next.
//	error: An error that will be non-nil if there was an issue querying the database.
func (crawled *CrawledUrl) GetNextCrawlUrls(limit int) ([]CrawledUrl, error) {
	var urls []CrawledUrl
	tx := DBConn.Where("last_tested IS NILL").Limit(limit).Find(&urls)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return []CrawledUrl{}, tx.Error
	}
	return urls, nil
}

// Save persists the CrawledUrl instance to the database.
// It uses the global DBConn to save the instance.
// If there's an error during the save operation, it prints the error and returns it.
//
// Returns:
//
//	error: An error that will be non-nil if there was an issue saving the instance.
func (crawled *CrawledUrl) Save() error {
	tx := DBConn.Save(&crawled)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return tx.Error
	}
	return nil
}

// GetNotIndex retrieves all CrawledUrl instances that have not been indexed yet.
// It uses the global DBConn to query the database.
// It returns a slice of CrawledUrl instances and an error.
// If there's an error during the query, it returns an empty slice and the error.
//
// Returns:
//
//	[]CrawledUrl: A slice of CrawledUrl instances that have not been indexed yet.
//	error: An error that will be non-nil if there was an issue querying the database.
func (crawled *CrawledUrl) GetNotIndex() ([]CrawledUrl, error) {
	var urls []CrawledUrl
	tx := DBConn.Where("indexed = > aAND last_tested IS NOT NULL", false).Find(&urls)
	if tx.Error != nil {
		return []CrawledUrl{}, tx.Error
	}
	return urls, nil
}

// SetIndexedTrue sets the Indexed field to true for each CrawledUrl instance in the given slice.
// It uses the global DBConn to save each instance to the database.
// If there's an error during the save operation, it returns the error immediately.
//
// Parameters:
//
//	urls: A slice of CrawledUrl instances to update.
//
// Returns:
//
//	error: An error that will be non-nil if there was an issue saving an instance.
func (crawled *CrawledUrl) SetIndexedTrue(urls []CrawledUrl) error {
	for _, url := range urls {
		url.Indexed = true
		tx := DBConn.Save(&url)

		if tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}

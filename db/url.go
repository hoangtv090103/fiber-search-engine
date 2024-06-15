package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CrawledUrl struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Url             string         `json:"url" gorm:"unique;not null"`
	Success         bool           `json:"success" gorm:"default:null"`
	CrawlDuration   time.Duration  `json:"crawlDuration"`
	ResponseCode    int            `json:"responseCode" gorm:"type:smallint"`
	PageTitle       string         `json:"pageTitle"`
	PageDescription string         `json:"pageDescription"`
	Headings        string         `json:"headings"`
	LastTested      *time.Time     `json:"lastTested"` // Use pointer so this value can be nil
	Indexed         bool           `json:"indexed" gorm:"default:false"`
	CreatedAt       *time.Time     `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

// GetUrl is a method on the CrawledUrl struct that retrieves a crawled URL from the database.
// It fetches the crawled URL with the specified ID and populates the CrawledUrl struct with the retrieved values.
//
// Parameters:
// id string: The ID of the crawled URL to retrieve.
//
// Returns:
// error: An error object that describes an error that occurred during the method's execution.
func (crawled *CrawledUrl) UpdateUrl(input CrawledUrl) error {
	tx := DBConn.Select("url", "success", "crawl_duration", "response_code", "page_title", "page_description", "headings", "last_tested", "updated_at").Omit("created_at").Save(&input)
	if tx.Error != nil {
		fmt.Print(tx.Error)
		return tx.Error
	}
	return nil
}

// GetNextCrawlUrls is a method on the CrawledUrl struct that retrieves the next URLs to crawl from the database.
// It fetches the URLs that have not been tested yet and limits the number of URLs to the specified limit.
//
// Parameters:
// limit int: The maximum number of URLs to retrieve.
//
// Returns:
// []CrawledUrl: A slice of CrawledUrl objects representing the URLs to crawl.
// error: An error object that describes an error that occurred during the method's execution.
func (crawled *CrawledUrl) GetNextCrawlUrls(limit int) ([]CrawledUrl, error) {
	var urls []CrawledUrl
	tx := DBConn.Where("last_tested IS NULL").Limit(limit).Find(&urls)
	if tx.Error != nil {
		fmt.Print(tx.Error)
		return []CrawledUrl{}, tx.Error
	}
	return urls, nil
}

// Save is a method on the CrawledUrl struct that saves the crawled URL to the database.
// It takes a CrawledUrl object that represents the URL that has been crawled and saves it to the database.
//
// This method does not take any parameters.
//
// Returns:
// error: An error object that describes an error that occurred during the method's execution.
func (crawled *CrawledUrl) Save() error {
	tx := DBConn.Save(&crawled)
	if tx.Error != nil {
		fmt.Print(tx.Error)
		return tx.Error
	}
	return nil
}

// GetNotIndexed is a method on the CrawledUrl struct that retrieves the URLs that have not been indexed from the database.
// It fetches the URLs that have not been indexed and have been tested, and returns them as a slice of CrawledUrl objects.
//
// This method does not take any parameters.
//
// Returns:
// []CrawledUrl: A slice of CrawledUrl objects representing the URLs that have not been indexed.
// error: An error object that describes an error that occurred during the method's execution.
func (crawled *CrawledUrl) GetNotIndexed() ([]CrawledUrl, error) {
	var urls []CrawledUrl
	tx := DBConn.Where("indexed = ? AND last_tested IS NOT NULL", false).Find(&urls)
	if tx.Error != nil {
		fmt.Print(tx.Error)
		return []CrawledUrl{}, tx.Error
	}
	return urls, nil
}

// SetIndexedTrue is a method on the CrawledUrl struct that sets the indexed flag to true for a slice of CrawledUrl objects.
// It takes a slice of CrawledUrl objects and sets the indexed flag to true for each URL in the slice.
//
// Parameters:
// urls []CrawledUrl: A slice of CrawledUrl objects representing the URLs to mark as indexed.
//
// Returns:
// error: An error object that describes an error that occurred during the method's execution.
func (crawled *CrawledUrl) SetIndexedTrue(urls []CrawledUrl) error {
	for _, url := range urls {
		url.Indexed = true
		tx := DBConn.Save(&url)
		if tx.Error != nil {
			fmt.Print(tx.Error)
			return tx.Error
		}
	}
	return nil
}

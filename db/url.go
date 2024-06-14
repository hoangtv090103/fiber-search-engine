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

func (crawled *CrawledUrl) UpdatedUrl(input CrawledUrl) error {
	tx := DBConn.Select("url", "success", "crawl_duration", "response_code", "page_title", "page_description", "headings", "last_tested", "updated_at").Omit("create_ai").Save(&input)

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return tx.Error
	}
	return nil
}

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

package search

import (
	"fiber-search-engine/db"
	"fmt"
	"time"
)

func RunEngine() {
	fmt.Println("Started search engine crawl...")
	defer fmt.Println("Finished search engine crawl")

	// Get settings
	settings := &db.SearchSetting{}
	err := settings.Get()
	if err != nil {
		fmt.Println("something went wrong getting settings")
		return
	}
	// Check if search is on
	if !settings.SearchOn {
		fmt.Println("Search is turned off")
		return
	}

	crawl := &db.CrawledUrl{}
	nextUrls, err := crawl.GetNextCrawlUrls(int(settings.Amount))

	if err != nil {
		fmt.Println("something went wrong getting next urls")
		return
	}
	newUrls := []db.CrawledUrl{}

	testedTime := time.Now()

	for _, next := range nextUrls {
		result := runCrawl(next.Url)
		if !result.Success {
			err := next.UpdatedUrl(db.CrawledUrl{
				ID:              next.ID,
				Url:             next.Url,
				Success:         false,
				CrawlDuration:   result.CrawlData.CrawlTime,
				ResponseCode:    result.ResponseCode,
				PageTitle:       crawl.PageTitle,
				PageDescription: result.CrawlData.PageDescription,
				Heading:         result.CrawlData.Headings,
			})

			if err != nil {
				fmt.Println("something went wrong updating a failed url")
				return
			}
			continue
		}

		// Success
		err := next.UpdatedUrl(db.CrawledUrl{
			ID:              next.ID,
			Url:             next.Url,
			Success:         result.Success,
			CrawlDuration:   result.CrawlData.CrawlTime,
			ResponseCode:    result.ResponseCode,
			PageTitle:       result.CrawlData.PageTitle,
			PageDescription: result.CrawlData.PageDescription,
			Heading:         result.CrawlData.Headings,
			LastTested:      &testedTime,
		})
		if err != nil {
			fmt.Println("something went wrong updating a successful url")
			fmt.Println(next.Url)
		}

		for _, newUrl := range result.CrawlData.Links.Internal {
			newUrls = append(newUrls, db.CrawledUrl{
				Url: newUrl,
			})
		}
	} // end of range

	if !settings.AddNew {
		return
	}

	// Insert new urls
	for _, newUrl := range newUrls {
		err := newUrl.Save()
		if err != nil {
			fmt.Println("something went wrong saving a new url")
			fmt.Println(newUrl.Url)
		}
	}

	fmt.Printf("\nAdded %d new urls to the database\n", len(newUrls))
}

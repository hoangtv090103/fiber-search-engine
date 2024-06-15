package search

import (
	"fiber-search-engine/db"
	"fmt"
	"time"
)

// RunEngine starts the search engine crawl process.
// It first retrieves the search settings and checks if the search is turned on.
// If the search is off, it returns immediately.
// It then retrieves the next URLs to crawl based on the amount specified in the settings.
// For each URL, it runs a crawl and updates the URL in the database with the crawl result.
// If the crawl is successful, it also adds the internal links found during the crawl to the newUrls slice.
// If the AddNew setting is true, it saves all new URLs in the newUrls slice to the database.
// It prints the number of new URLs added to the database at the end.
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

// RunIndex starts the search engine indexing process.
// It first retrieves all URLs that have not been indexed yet.
// It then creates a new index and adds the not indexed URLs to it.
// After that, it saves the index to the database.
// Finally, it sets the indexed field to true for all the not indexed URLs.
// If there's an error during any of these operations, it returns immediately.
func RunIndex() {
	fmt.Println("Started search engine index...")
	defer fmt.Println("Finished search engine index")

	crawled := &db.CrawledUrl{}              // create a new instance of the db.CrawledUrl struct
	notIndexed, err := crawled.GetNotIndex() // get all the not indexed urls
	if err != nil {
		return
	}

	idx := make(Index)               // create a new instance of the Index map
	idx.Add(notIndexed)              // add the notIndexed slice to the index
	searchIndex := &db.SearchIndex{} // create a new instance of the db.SearchIndex struct

	err = searchIndex.Save(idx, notIndexed) // save the index to the database
	if err != nil {
		return
	}

	err = crawled.SetIndexedTrue(notIndexed) // set the indexed field to true for all the not indexed urls
	if err != nil {
		return
	}
}

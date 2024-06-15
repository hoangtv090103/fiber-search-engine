package search

import (
	"fiber-search-engine/db"
	"fmt"
	"time"
)

// RunEngine is a function that runs the web crawling engine.
// It first prints a message that the crawl has started and defers a message that the crawl has finished.
// It then retrieves the crawl settings from the database and checks if search is turned on.
// If there is an error retrieving the settings or if search is turned off, it prints a message and returns.
// The function then retrieves the next set of URLs to be crawled from the database.
// If there is an error retrieving the URLs, it prints a message and returns.
// The function then loops over the URLs, runs a crawl on each one, and updates the database with the results.
// If the crawl is not successful, it updates the database with the failed crawl and continues to the next URL.
// If the crawl is successful, it updates the database with the successful crawl and adds the newly found external URLs to a slice.
// After all URLs have been crawled, the function checks if it should add the newly found URLs to the database.
// If it should, it loops over the new URLs and adds each one to the database.
// If there is an error adding a URL to the database, it prints a message.
// Finally, the function prints a message with the number of new URLs added to the database.
//
// This function does not take any parameters and does not return any values.
func RunEngine() {
	fmt.Println("started search engine crawl...")
	defer fmt.Println("search engine crawl has finished")
	// Get crawl settings from DB
	settings := &db.SearchSettings{}
	err := settings.Get()
	if err != nil {
		fmt.Println("something went wrong getting the settings")
		return
	}
	// Check if search is turned on by checking settings
	if !settings.SearchOn {
		fmt.Println("search is turned off")
		return
	}
	crawl := &db.CrawledUrl{}
	// Get next X urls to be tested
	nextUrls, err := crawl.GetNextCrawlUrls(int(settings.Amount))
	if err != nil {
		fmt.Println("something went wrong getting the url list")
		return
	}
	newUrls := []db.CrawledUrl{}
	testedTime := time.Now()
	// Loop over the slice and run crawl on each url
	for _, next := range nextUrls {
		result := runCrawl(next.Url)
		// Check if the crawl was not successul
		if !result.Success {
			// Update row in database with the failed crawl
			err := next.UpdateUrl(db.CrawledUrl{
				ID:              next.ID,
				Url:             next.Url,
				Success:         false,
				CrawlDuration:   result.CrawlData.CrawlTime,
				ResponseCode:    result.ResponseCode,
				PageTitle:       result.CrawlData.PageTitle,
				PageDescription: result.CrawlData.PageDescription,
				Headings:        result.CrawlData.Headings,
				LastTested:      &testedTime,
			})
			if err != nil {
				fmt.Println("something went wrong updating a failed url")
			}
			continue
		}
		// Update a successful row in database
		err := next.UpdateUrl(db.CrawledUrl{
			ID:              next.ID,
			Url:             next.Url,
			Success:         result.Success,
			CrawlDuration:   result.CrawlData.CrawlTime,
			ResponseCode:    result.ResponseCode,
			PageTitle:       result.CrawlData.PageTitle,
			PageDescription: result.CrawlData.PageDescription,
			Headings:        result.CrawlData.Headings,
			LastTested:      &testedTime,
		})
		if err != nil {
			fmt.Printf("something went wrong updating %v /n", next.Url)
		}
		// Push the newly found external urls to an array
		for _, newUrl := range result.CrawlData.Links.External {
			newUrls = append(newUrls, db.CrawledUrl{Url: newUrl})
		}
	} // End of range
	// Check if we should add the newly found urls to the database
	if !settings.AddNew {
		fmt.Printf("Adding new urls to database is disabled")
		return
	}
	// Insert newly found urls into database
	for _, newUrl := range newUrls {
		err := newUrl.Save()
		if err != nil {
			fmt.Printf("something went wrong adding new url to database: %v", newUrl.Url)
		}
	}
	fmt.Printf("\nAdded %d new urls to database \n", len(newUrls))
}

// RunIndex is a function that runs the search indexing process.
// It first prints a message that the indexing has started and defers a message that the indexing has finished.
// It then retrieves all URLs that have not been indexed from the database.
// If there is an error retrieving the URLs, it prints a message and returns.
// The function then creates a new index and adds the not indexed URLs to it.
// It then saves the index to the database.
// If there is an error saving the index, it prints a message and returns.
// Finally, it updates the URLs in the database to be indexed=true.
// If there is an error updating the URLs, it prints a message and returns.
//
// This function does not take any parameters and does not return any values.
func RunIndex() {
	fmt.Println("started search indexing...")
	defer fmt.Println("search indexing has finished")
	// Get index settings from DB
	crawled := &db.CrawledUrl{}
	// Get all urls that are not indexed
	notIndexed, err := crawled.GetNotIndexed()
	fmt.Println("not indexed urls: ", len(notIndexed))
	if err != nil {
		fmt.Println("something went wrong getting the not indexed urls")
		return
	}
	// Create a new index
	idx := make(Index)
	// Add the not indexed urls to the index
	idx.Add(notIndexed)
	// Save the index to the database
	searchIndex := &db.SearchIndex{}
	err = searchIndex.Save(idx, notIndexed)
	if err != nil {
		fmt.Println(err)
		fmt.Println("something went wrong saving the index")
		return
	}
	// Update the urls to be indexed=true
	err = crawled.SetIndexedTrue(notIndexed)
	if err != nil {
		fmt.Println("something went wrong updating the indexed urls")
		return
	}

}

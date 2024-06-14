package utils

import (
	"fiber-search-engine/search"
	"fmt"

	"github.com/robfig/cron"
)

// StartCronJobs initializes and starts the cron jobs for the application.
// It creates a new cron instance, adds the search engine run function to run every hour, and starts the cron jobs.
// It also prints the number of cron jobs that have been set up.
func StartCronJobs() {
	c := cron.New()
	// add cron jobs here
	c.AddFunc("@every 1h", search.RunEngine)
	c.Start()
	croneCount := len(c.Entries())
	fmt.Printf("setup %d cron jobs\n", croneCount)
}

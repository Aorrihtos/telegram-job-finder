package scrapper

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var scrappers = []Scrapper{
	&scrapeRemoteOK{},
	&ScrapeWeWorkRemotely{},
}

func RunScrapper() {
	httpClient := &http.Client{
		Timeout: 10 * time.Minute,
	}

	ticker := time.NewTicker(10 * time.Minute)

	// Launch at startup and then, every minute
	wg := &sync.WaitGroup{}
	for ; ; <-ticker.C {
		launchScrappers(wg, httpClient)
	}
}

func launchScrappers(wg *sync.WaitGroup, httpClient *http.Client) {
	insertedJobs := make([]any, 0, 4000) // Default capacity for storage
	for _, s := range scrappers {
		wg.Add(1)
		go s.scrape(httpClient, wg, &insertedJobs)
	}

	wg.Wait()
	// TODO: Notify the users with the new elements
	log.Printf("Total inserted jobs this run: %v", len(insertedJobs))
}
package scrapper

import (
	"net/http"
	"sync"
	"time"
)

var scrappers = []Scrapper{
	&scrapeRemoteOK{},
}

func RunScrapper() {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	ticker := time.NewTicker(1 * time.Minute)

	// Launch at startup and then, every minute
	wg := &sync.WaitGroup{}
	for ; ; <-ticker.C {
		launchScrappers(wg, httpClient)
	}
}

func launchScrappers(wg *sync.WaitGroup, httpClient *http.Client) {
	for _, s := range scrappers {
		wg.Add(1)
		go s.scrape(httpClient, wg)
	}

	wg.Wait()
}
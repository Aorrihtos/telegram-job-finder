package scrapper

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Scrapper interface {
	scrape(*http.Client,*sync.WaitGroup)
}

type scrapeRemoteOK struct {}

func (s *scrapeRemoteOK) scrape(httpClient *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := httpClient.Get("https://remoteok.com/api")
	if err != nil {
		log.Println("Error fetching jobs from remoteOk:", err)
		return
	}

	defer resp.Body.Close()

	// Read and unmarshal the response...
	var jobs []RemoteOKJob
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		log.Println("Error decoding JSON response:", err)
		return
	}

	// Process the response...
	log.Println("Fetched jobs:", jobs)
}
package scrapper

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Aorrihtos/telegram-job-finder/db"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	var jobs []RemoteOKJob
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		log.Println("Error decoding JSON response:", err)
		return
	}

	jobs = jobs[1:] // Remove the first element which is not a job

	dbJobs := make([]Job, len(jobs))
	for i, job := range jobs {
		publishedDate, err := time.Parse("2006-01-02T15:04:05-07:00", job.Date)
		if err != nil {
			log.Println("Error parsing published date:", err)
			continue
		}
		dbJobs[i] = Job{
			Position:    job.Position,
			JobType:    "Undefined",
			Category:   "Undefined",
			Salary: struct {
				Min int `json:"min"`
				Max int `json:"max"`
			}{
				Min: job.SalaryMin,
				Max: job.SalaryMax,
			},
			Location:    job.Location,
			URL:         job.URL,
			PublishedDate: publishedDate,
		}
	}

	opts := options.InsertMany().SetOrdered(false)
	elems, _ := db.GetJobsCollection().InsertMany(context.Background(), dbJobs, opts)
	log.Printf("Inserted %d jobs from remoteOK into MongoDB", len(dbJobs))
	
	// TODO: Notify the users with the new elemes
	log.Printf("Inserted elements: %v", elems.InsertedIDs)
}
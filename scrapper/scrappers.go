package scrapper

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Aorrihtos/telegram-job-finder/db"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Scrapper interface {
	scrape(*http.Client,*sync.WaitGroup, *[]any)
}

type scrapeRemoteOK struct {}

func (s *scrapeRemoteOK) scrape(httpClient *http.Client, wg *sync.WaitGroup, insertedJobs *[]any) {
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

		categories := getCategoriesFromTitle(job.Position)
		if len(categories) == 0 {
			categories = []string{"other"}
		}

		dbJobs[i] = Job{
			Position:    job.Position,
			Type:    "full_time",
			Company:   job.Company,
			Category:   categories[0],
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

	log.Printf("RemoteOK inserted elements: %v", len(elems.InsertedIDs))
	*insertedJobs = append(*insertedJobs, elems.InsertedIDs...)
}

type ScrapeWeWorkRemotely struct{}

func (s *ScrapeWeWorkRemotely) scrape(httpClient *http.Client, wg *sync.WaitGroup, insertedJobs *[]any) {
	defer wg.Done()

	resp, err := httpClient.Get("https://weworkremotely.com/remote-jobs.rss")
	if err != nil {
		log.Println("Error fetching jobs from We Work Remotely:", err)
		return
	}

	defer resp.Body.Close()

	var channel WeWorkRemotelyChannel
	errParse := xml.NewDecoder(resp.Body).Decode(&channel)
	if errParse != nil {
		log.Println("Error decoding XML response:", errParse)
		return
	}

	jobs := make([]Job, len(channel.Channel.Items))
	for i, job := range channel.Channel.Items {
		pubDate, err := time.Parse(time.RFC1123Z, job.PubDate)
		if err != nil {
			log.Println("Error parsing published date:", err)
			continue
		}

		categories := getCategoriesFromTitle(job.Title)
		if len(categories) == 0 {
			categories = []string{"other"}
		}

		jobType := strings.ReplaceAll(job.Type, "-", "_")

		jobs[i] = Job{
			Position:    job.Title,
			Type:       strings.ToLower(jobType),
			Category:   categories[0],
			Location:    job.Region,
			URL:         job.Link,
			PublishedDate: pubDate,
		}
	}

	// Insert jobs into the database
	opts := options.InsertMany().SetOrdered(false)
	elems, _ := db.GetJobsCollection().InsertMany(context.Background(), jobs, opts)

	log.Printf("WeWorkRemotely inserted elements: %v", len(elems.InsertedIDs))
	*insertedJobs = append(*insertedJobs, elems.InsertedIDs...)
}
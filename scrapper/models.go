package scrapper

import "time"

type RemoteOKJob struct {
	Date        string `json:"date"`
	Position    string `json:"position"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	SalaryMin   int    `json:"salary_min"`
	SalaryMax   int    `json:"salary_max"`
	URL         string `json:"url"`
}

type Job struct {
	Position string `json:"position"`
	JobType      string `json:"job_type"`
	Category     string `json:"category"`
	Salary       struct{
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"salary"`
	Location string `json:"location"`
	URL      string `json:"url"`
	PublishedDate time.Time `json:"published_date"`
}
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

type WeWorkRemotelyJob struct {
	Title      string `xml:"title"`
	Region     string `xml:"region"`
	Country    string `xml:"country"`
	Skills     string `xml:"skills"`
	Category   string `xml:"category"`
	Type       string `xml:"type"`
	PubDate    string `xml:"pubDate"`
	Link       string `xml:"link"`
}

type WeWorkRemotelyChannel struct {
	Channel struct {
		Title string `xml:"title"`
		Link string `xml:"link"`
		Description string `xml:"description"`
		Language string `xml:"language"`
		TTL int `xml:"ttl"`
		Items []WeWorkRemotelyJob `xml:"item"`
	} `xml:"channel"`
}
type Job struct {
	Position string `json:"position"`
	Company string `json:"company"`
	Type      string `json:"type"`
	Category     string `json:"category"`
	Salary       struct{
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"salary"`
	Location string `json:"location"`
	URL      string `json:"url"`
	PublishedDate time.Time `json:"published_date"`
}
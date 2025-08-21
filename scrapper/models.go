package scrapper

type RemoteOKJob struct {
	Date        string `json:"date"`
	Position    string `json:"position"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	SalaryMin   int    `json:"salary_min"`
	SalaryMax   int    `json:"salary_max"`
	URL         string `json:"url"`
}

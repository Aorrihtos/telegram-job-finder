package telegram_bot

type SearchForm struct {
	JobType string `json:"job_type"`
	Specialization string `json:"specialization"`
	SalaryRange string `json:"salary_range"`
}
package telegram_bot

import (
	"strconv"
	"strings"
)

type Salary struct {
    Min int `json:"min"`
    Max int `json:"max"`
}

type Preferences struct {
    JobType  string `json:"job_type"`
    Category string `json:"category"`
    Salary   Salary `json:"salary"`
}

type SearchForm struct {
	JobType string `json:"job_type"`
	Category string `json:"category"`
	SalaryRange string `json:"salary_range"`
}

type User struct {
	ID int64 `json:"_id"`
	Preferences Preferences `json:"preferences"`
}

func (sf *SearchForm) toUserModel(chatId int64) User {
	minSalary, _ := strconv.Atoi(strings.Split(sf.SalaryRange, "-")[0])
	maxSalary, _ := strconv.Atoi(strings.Split(sf.SalaryRange, "-")[1])
	return User{
		ID: chatId,
		Preferences: Preferences{
			JobType: sf.JobType,
			Category: sf.Category,
			Salary: Salary{
				Min: minSalary,
				Max: maxSalary,
			},
		},
	}
}
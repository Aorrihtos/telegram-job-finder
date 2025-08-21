package telegram_bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var searchForm SearchForm

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if (update.Message != nil){
		switch update.Message.Text {
			case "/start":
				NewSessionHandler(ctx, b, update)
		}
	}

	if update.CallbackQuery != nil {
		HandleAnswer(ctx, b, update)
	}
}

// Handle user answer when button is clicked
func HandleAnswer(ctx context.Context, b *bot.Bot, update *models.Update) {
	cq := update.CallbackQuery
	answer := strings.Split(cq.Data, ":")
	if len(answer) < 2 { 
		fmt.Println("Answer does not compliment the format") 
		return 
	}

    switch answer[0] {
    case "job_type":
        searchForm.JobType = answer[1]
		AskSpecializationHandler(ctx, b, update)
    case "specialization":
        searchForm.Specialization = answer[1]
		AskSalaryRangeHandler(ctx, b, update)
    case "salary_range":
        searchForm.SalaryRange = answer[1]
		fetchDbJobs(ctx, b, update)
    }

}

func NewSessionHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Welcome to the job finder bot! Tell me your preferences and I will dig the internet for you",
	})

	msg := &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:  "Please select your preferred job type:",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text: "Full-time",
						CallbackData: "job_type:full_time",
					},
					{
						Text: "Part-time",
						CallbackData: "job_type:part_time",
					},
				},
			},
		},
	}

	b.SendMessage(ctx, msg)
}

func AskSpecializationHandler(ctx context.Context, b *bot.Bot, update *models.Update,){
	msg := &bot.SendMessageParams{
		ChatID: update.CallbackQuery.From.ID,
		Text:  "Please select your preferred specialization:",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text: "Full-Stack Developer",
						CallbackData: "specialization:fullstack",
					},
					{
						Text: "Back-end Developer",
						CallbackData: "specialization:backend",
					},
					{
						Text: "Front-end Developer",
						CallbackData: "specialization:frontend",
					},
					{
						Text: "Dev-Ops",
						CallbackData: "specialization:devops",
					},
				},
			},
		},
	}

	b.SendMessage(ctx, msg)
}

func AskSalaryRangeHandler(ctx context.Context, b *bot.Bot, update *models.Update){
	msg := &bot.SendMessageParams{
		ChatID: update.CallbackQuery.From.ID,
		Text:  "Please select your preferred salary range:",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text: "<= 20.000 euros",
						CallbackData: "salary_range:<=20.000",
					},
					{
						Text: "> 20.000 <= 30.000",
						CallbackData: "salary_range:>20.000<=30.000",
					},
					{
						Text: "> 30.000 <= 50.000",
						CallbackData: "salary_range:>30.000<=50.000",
					},
					{
						Text: "> 50.000",
						CallbackData: "salary_range:>50.000",
					},
				},
			},
		},
	}

	b.SendMessage(ctx, msg)
}

func fetchDbJobs(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg := fmt.Sprintf("This were your preferences: \nJob Type: %s\nSpecialization: %s\nSalary Range: %s",
		searchForm.JobType, searchForm.Specialization, searchForm.SalaryRange)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.From.ID,
		Text:   msg,
	})
}
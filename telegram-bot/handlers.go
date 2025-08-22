package telegram_bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/Aorrihtos/telegram-job-finder/db"
	"github.com/Aorrihtos/telegram-job-finder/scrapper"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
    case "category":
        searchForm.Category = answer[1]
		AskSalaryRangeHandler(ctx, b, update)
    case "salary_range":
        searchForm.SalaryRange = answer[1]
		saveUserPreferences(ctx, b, update)
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
					{
						Text: "Contract",
						CallbackData: "job_type:contract",
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
						CallbackData: "category:fullstack",
					},
					{
						Text: "Back-end Developer",
						CallbackData: "category:backend",
					},
					{
						Text: "Front-end Developer",
						CallbackData: "category:frontend",
					},
					{
						Text: "Mobile Developer",
						CallbackData: "category:mobile",
					},
					{
						Text: "Data Scientist",
						CallbackData: "category:data",
					},
					{
						Text: "AI Developer",
						CallbackData: "category:ml_ai",
					},
					{
						Text: "Dev-Ops",
						CallbackData: "category:devops",
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
						CallbackData: "salary_range:0-20000",
					},
					{
						Text: "> 20.000 <= 30.000",
						CallbackData: "salary_range:20000-30000",
					},
					{
						Text: "> 30.000 <= 50.000",
						CallbackData: "salary_range:30000-50000",
					},
					{
						Text: "> 50.000",
						CallbackData: "salary_range:50000-999999",
					},
				},
			},
		},
	}

	b.SendMessage(ctx, msg)
}

func saveUserPreferences(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg := fmt.Sprintf("This were your preferences: \nJob Type: %s\nSpecialization: %s\nSalary Range: %s",
		searchForm.JobType, searchForm.Category, searchForm.SalaryRange)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.From.ID,
		Text:   msg,
	})

	// Save the user preferences in the db
	user := searchForm.toUserModel(update.CallbackQuery.From.ID)
	opts := options.UpdateOne().SetUpsert(true)
	_, err := db.GetUsersCollection().UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user.Preferences}, opts)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Error saving your preferences. Please try again later.",
		})
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.From.ID,
		Text:   "Your preferences have been saved successfully!",
	})

	sendUserAvailableJobs(ctx, b, update, user)
}

func sendUserAvailableJobs(ctx context.Context, b *bot.Bot, update *models.Update, user User) {
	currentJobsPosted, err := db.GetJobsCollection().Find(ctx, bson.M{
		"category":  user.Preferences.Category,
		"type":  user.Preferences.JobType,
		"salary.min":    bson.M{"$gte": user.Preferences.Salary.Min},
		"salary.max":    bson.M{"$lte": user.Preferences.Salary.Max},
	})
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Error fetching jobs from the database. Please try again later.",
		})
		return
	}

	defer currentJobsPosted.Close(ctx)

	var jobs []scrapper.Job
	if err := currentJobsPosted.All(ctx, &jobs); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Error fetching jobs from the database. Please try again later.",
		})
		return
	}

	if len(jobs) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "No jobs found matching your preferences at this time. When a new job is posted that matches your preferences, you'll be the first to know!",
		})
		return
	}

	for _, job := range jobs {
		jobMessage := fmt.Sprintf("Position: %s\nLocation: %s\nSalary: %d - %d euros\nURL: %s",
			job.Position, job.Location, job.Salary.Min, job.Salary.Max, job.URL)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   jobMessage,
		})
	}
}
package main

import (
	"context"

	"log"
	"os"
	"os/signal"

	"github.com/Aorrihtos/telegram-job-finder/db"
	"github.com/Aorrihtos/telegram-job-finder/scrapper"
	botFinderModels "github.com/Aorrihtos/telegram-job-finder/telegram-bot"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Create Telegram BOT
	botToken := os.Getenv("BOT_TOKEN")
	opts := []bot.Option{
		bot.WithDefaultHandler(botFinderModels.DefaultHandler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		log.Fatal("Error creating bot: ", err)
	}

	// Connect to MongoDB
	mongoClient, err := db.ConnectToDb()
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalln("Error disconnecting from MongoDB: ", err)
		}
	}()

	// Start Scrapper and Bot
	go scrapper.RunScrapper()
	go b.Start(ctx)

	log.Println("Bot started successfully. Listening for updates...")
	// Wait for an interrupt signal to gracefully shut down the bot
	<-ctx.Done()

	log.Println("Shutting down bot...")
}

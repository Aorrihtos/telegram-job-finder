package main

import (
	"context"
	"fmt"

	"log"
	"os"
	"os/signal"

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
	botToken := os.Getenv("BOT_TOKEN")

	opts := []bot.Option{
		bot.WithDefaultHandler(botFinderModels.DefaultHandler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		log.Fatal("Error creating bot: ", err)
	}

	go scrapper.RunScrapper()
	go b.Start(ctx)
	
	fmt.Println("Bot started successfully. Listening for updates...")
	// Wait for an interrupt signal to gracefully shut down the bot
	<-ctx.Done()

	fmt.Println("Shutting down bot...")
}

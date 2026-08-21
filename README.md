# telegram-job-finder

A Telegram bot that helps you find remote developer jobs. Message the bot on Telegram, tell it your job type, specialization, and salary range through a few button taps, and it replies with matching remote job listings. New listings are scraped from remote job boards every 10 minutes and matched against your saved preferences.

## Usage

1. Start a chat with the bot and send `/start`.
2. Pick your preferred job type (Full-time, Part-time, Contract).
3. Pick your specialization (Full-Stack, Back-end, Front-end, Mobile, Data, AI, DevOps).
4. Pick your salary range.
5. The bot saves your preferences and immediately sends any currently matching job listings.

## Job sites scraped

- [RemoteOK](https://remoteok.com)
- [We Work Remotely](https://weworkremotely.com)

## Requirements

- Go 1.25+
- A MongoDB instance
- A Telegram bot token ([create one via BotFather](https://core.telegram.org/bots#botfather))

## Configuration

Create a `.env` file in the project root with:

```
BOT_TOKEN=your-telegram-bot-token
MONGO_URL=your-mongodb-connection-string
```

## Running

```
go run main.go
```

This connects to MongoDB, starts the Telegram bot, and launches the job scrapers (running immediately, then every 10 minutes) in the background.
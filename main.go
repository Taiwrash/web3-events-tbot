package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Taiwrash/web3event-spot/scrape"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		fmt.Println("can not create bot!")
		panic(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		if update.Message == nil {
			continue
		}
		evts, _ := scrape.Scrape()
		if len(evts) == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No event found for the provided date.")
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		} else {
			for i, event := range evts {
				// Split the Date field into date and location parts
				parts := strings.SplitN(event.Date, "AM", 2)
				if len(parts) < 2 {
					parts = strings.SplitN(event.Date, "PM", 2)
				}
				if len(parts) == 2 {
					event.Date = strings.TrimSpace(parts[0]) + "AM"
					event.Location = strings.TrimSpace(parts[1])
				}

				msgText := fmt.Sprintf("Event %d:\nEvent Name: %s\nDate: %s\nVenue: %s\nRegistration Link: %s",
					i+1, event.Title, event.Date, event.Location, event.URL)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ReplyToMessageID = update.Message.MessageID
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}
		}
	}

	bot.Debug = true
}

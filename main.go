package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/louispy/linebot/internal/app"
	"github.com/louispy/linebot/internal/usecases/command"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := command.NewLineCommand(command.CommandOpts{
		Bot: bot,
	})
	cmd.LoadHandlers()

	commandApp := app.NewCommandApp(app.CommandAppOpts{
		Command: cmd,
	})

	commandApp.Run()
}

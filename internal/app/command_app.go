package app

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	cmd "github.com/louispy/linebot/internal/adapters"
)

type CommandApp struct {
	command cmd.Command
}

type CommandAppOpts struct {
	Command cmd.Command
}

func NewCommandApp(o CommandAppOpts) CommandApp {
	return CommandApp{
		command: o.Command,
	}
}

func (app CommandApp) Run() {
	log.Println("Starting Command App...")
	lambda.StartWithOptions(app.command.Callback, lambda.WithContext(context.Background()))
}

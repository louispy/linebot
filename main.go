package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

type LineWebhookRequest struct {
	Destination string           `json:"destination"`
	Events      []*linebot.Event `json:"events"`
}

func callback(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var lineReq LineWebhookRequest

	err := json.Unmarshal([]byte(req.Body), &lineReq)
	if err != nil {
		log.Println(err)
		return &events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	for _, event := range lineReq.Events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}

	return &events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}
	bot, err = linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	lambda.Start(callback)
}

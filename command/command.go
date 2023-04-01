package command

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/line/line-bot-sdk-go/linebot"
)

type LineWebhookRequest struct {
	Destination string           `json:"destination"`
	Events      []*linebot.Event `json:"events"`
}

type Command struct {
	bot *linebot.Client
}

type CommandOpts struct {
	Bot *linebot.Client
}

func NewCommand(o CommandOpts) Command {
	return Command{
		bot: o.Bot,
	}
}

func (c Command) Callback(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
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
				msg := message.Text
				if len(msg) > 0 && msg[0] == '/' {
					args := strings.Split(msg, " ")
					cmd := args[0][1:]
					args = args[1:]
					switch cmd {
					case "search":
						if len(args) > 0 {
							c.Search(event, strings.Join(args, " "))
						}
					default:
						if _, err := c.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Unrecognized Command :(")).Do(); err != nil {
							log.Println(err)
						}

					}
				}
			}
		}
	}

	return &events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

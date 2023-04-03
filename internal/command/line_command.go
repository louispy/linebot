package command

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/louispy/linebot/internal/constants"
	"go.opentelemetry.io/otel"
)

type LineWebhookRequest struct {
	Destination string           `json:"destination"`
	Events      []*linebot.Event `json:"events"`
}

type LineCommand struct {
	bot *linebot.Client
}

type CommandOpts struct {
	Bot *linebot.Client
}

func NewLineCommand(o CommandOpts) LineCommand {
	return LineCommand{
		bot: o.Bot,
	}
}

func (c LineCommand) Callback(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ctx, span := otel.Tracer(constants.Usecase).Start(ctx, constants.Callback)
	defer span.End()

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
							c.Search(ctx, event, strings.Join(args, " "))
						}
					default:
						if _, err := c.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Unrecognized Command :(")).WithContext(ctx).Do(); err != nil {
							log.Println(err)
						}

					}
				}
			}
		}
	}

	return &events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

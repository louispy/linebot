package adapters

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Command interface {
	Callback(context.Context, *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
	Search(ctx context.Context, event *linebot.Event, query string)
}

package adapters

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type Command interface {
	Callback(context.Context, *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
	LoadHandlers()
}

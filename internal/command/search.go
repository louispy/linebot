package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/louispy/linebot/internal/constants"
	"go.opentelemetry.io/otel"
)

type searchResult struct {
	Heading string `json:"Heading"`
	Text    string `json:"AbstractText"`
	Url     string `json:"AbstractURL"`
}

func (c Command) Search(ctx context.Context, event *linebot.Event, query string) {
	ctx, span := otel.Tracer(constants.Usecase).Start(ctx, constants.CommandSearch)
	defer span.End()

	url := fmt.Sprintf("%s&q=%s", os.Getenv("SEARCH_API_BASE_URL"), query)

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var result searchResult
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		log.Println(err)
		return
	}

	message := "Sorry! No search result..."

	if result.Heading != "" {
		message = result.Heading
		if result.Text != "" {
			text := result.Text
			if len(text) > 100 {
				text = text[0:100] + "..."
			}
			message += fmt.Sprintf("\n%s", text)
		}
		if result.Url != "" {
			message += fmt.Sprintf("\n%s", result.Url)
		}
	}

	_, err = c.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).WithContext(ctx).Do()
	if err != nil {
		log.Println(err)
	}
}

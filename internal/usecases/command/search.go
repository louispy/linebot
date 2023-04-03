package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/louispy/linebot/internal/constants"
	"go.opentelemetry.io/otel"
)

type searchResult struct {
	Heading string `json:"title"`
	Text    string `json:"description"`
	Url     string `json:"url"`
}

type searchResponse struct {
	Query   string         `json:"query"`
	Results []searchResult `json:"results"`
}

func (c LineCommand) Search(ctx context.Context, event *linebot.Event, query string) {
	ctx, span := otel.Tracer(constants.Usecase).Start(ctx, constants.CommandSearch)
	defer span.End()

	host := os.Getenv("SEARCH_API_BASE_URL")
	apiKey := os.Getenv("RAPID_API_KEY")

	url := fmt.Sprintf("https://%s?q=%s", host, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Add("X-RapidAPI-Key", apiKey)
	req.Header.Add("X-RapidAPI-Host", host)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var response searchResponse
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		log.Println(err)
		return
	}

	message := "Sorry! No search result..."

	if len(response.Results) > 0 {
		result := response.Results[0]
		message = result.Heading
		if result.Text != "" {
			text := result.Text
			if len(text) > 150 {
				text = text[0:150] + "..."
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

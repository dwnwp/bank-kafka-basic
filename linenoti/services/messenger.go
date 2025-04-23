package services

import (
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type App struct {
	channelSecret string
	bot           *messaging_api.MessagingApiAPI
}

func NewApp(channelSecret, channelToken string) (*App, error) {
	bot, err := messaging_api.NewMessagingApiAPI(channelToken)
	if err != nil {
		return &App{}, err
	}
	return &App{
		channelSecret: channelSecret,
		bot:           bot,
	}, err
}

func (app *App) HandleEvents(c *gin.Context) {
	req, err := webhook.ParseRequest(app.channelSecret, c.Request)
	if err != nil {
		log.Println("Error parsing request:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	log.Println("[Line messaging] handling events...")

	for _, event := range req.Events {
		log.Printf("[Line messaging] /linebot called %+v...\n", event)

		switch e := event.(type) {
		case webhook.MessageEvent:
			app.bot.ReplyMessage(
				&messaging_api.ReplyMessageRequest{
					ReplyToken: e.ReplyToken,
					Messages: []messaging_api.MessageInterface{
						messaging_api.TextMessage{
							Text: "Hello World",
						},
					},
				},
			)
		}
	}

	c.Status(http.StatusOK)
}

func (app *App) SendDepositNotify(message *sarama.ConsumerMessage) {
	app.bot.PushMessage(
		&messaging_api.PushMessageRequest{
			To: "",
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: message.Topic,
				},
			},
		},
		"", // x-line-retry-key
	)
}

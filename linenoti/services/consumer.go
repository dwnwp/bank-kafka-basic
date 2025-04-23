package services

import (
	"encoding/json"
	"go-kafka-bank/events"
	"log"
	"reflect"

	"github.com/IBM/sarama"
)

type Consumer struct {
	app *App
}

func NewConsumer(app *App) *Consumer {
	return &Consumer{
		app: app,
	}
}

func (consumer *Consumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }

func (consumer *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	
	for message := range claim.Messages() {

		switch message.Topic {
			case reflect.TypeOf(events.OpenAccountEvent{}).Name():
				event := &events.OpenAccountEvent{}
				err := json.Unmarshal(message.Value, event)
				if err != nil {
					log.Println(err)
					return err
				}

				log.Printf("[%v] %+v", message.Topic, event)

			case reflect.TypeOf(events.DepositFundEvent{}).Name():
				event := &events.DepositFundEvent{}
				err := json.Unmarshal(message.Value, event)
				if err != nil {
					log.Println(err)
					return err
				}
				consumer.app.SendDepositNotify(message)
				log.Printf("[%v] %+v", message.Topic, event)

			case reflect.TypeOf(events.WithdrawFundEvent{}).Name():
				event := &events.WithdrawFundEvent{}
				err := json.Unmarshal(message.Value, event)
				if err != nil {
					log.Println(err)
					return err
				}
				log.Printf("[%v] %+v", message.Topic, event)

			case reflect.TypeOf(events.CloseAccountEvent{}).Name():
				event := &events.CloseAccountEvent{}
				err := json.Unmarshal(message.Value, event)
				if err != nil {
					log.Println(err)
					return err
				}
				log.Printf("[%v] %+v", message.Topic, event)

			default:
				log.Println("no event handler")
		}

		session.MarkMessage(message, "")
	}
	return nil
}

package services

import (
	"encoding/json"
	"go-kafka-bank/consumer/models"
	"go-kafka-bank/events"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EventHandler interface {
	Handle(topic string, eventBytes []byte)
}

type accountEventHandler struct {
	accountRepo models.AccountRepository
}

func NewAccountEventHandler(accountRepo models.AccountRepository) EventHandler {
	return accountEventHandler{accountRepo}
}

func (obj accountEventHandler) Handle(topic string, eventBytes []byte) {
	switch topic {
	case reflect.TypeOf(events.OpenAccountEvent{}).Name():
		event := &events.OpenAccountEvent{}
		err := json.Unmarshal(eventBytes, event)
		if err != nil {
			log.Println(err)
			return
		}
		eventId, _ := bson.ObjectIDFromHex(event.ID)
		bankAccount := models.BankAccount{
			ID:            eventId,
			AccountHolder: event.AccountHolder,
			AccountType:   event.AccountType,
			Balance:       event.OpeningBalance,
		}
		err = obj.accountRepo.Save(bankAccount)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("[%v] %+v", topic, event)
	case reflect.TypeOf(events.DepositFundEvent{}).Name():
		event := &events.DepositFundEvent{}
		err := json.Unmarshal(eventBytes, event)
		if err != nil {
			log.Println(err)
			return
		}
		bankAccount, err := obj.accountRepo.FindByID(event.ID)
		if err != nil {
			log.Println(err)
			return
		}
		bankAccount.Balance += event.Amount

		err = obj.accountRepo.Update(bankAccount)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("[%v] %+v", topic, event)
	case reflect.TypeOf(events.WithdrawFundEvent{}).Name():
		event := &events.WithdrawFundEvent{}
		err := json.Unmarshal(eventBytes, event)
		if err != nil {
			log.Println(err)
			return
		}
		bankAccount, err := obj.accountRepo.FindByID(event.ID)
		if err != nil {
			log.Println(err)
			return
		}

		if event.Amount > bankAccount.Balance {
			log.Printf("[%v] Cannot withdraw balance < amount!", topic)
			return
		}
		bankAccount.Balance -= event.Amount

		err = obj.accountRepo.Update(bankAccount)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("[%v] %+v", topic, event)
	case reflect.TypeOf(events.CloseAccountEvent{}).Name():
		event := &events.CloseAccountEvent{}
		err := json.Unmarshal(eventBytes, event)
		if err != nil {
			log.Println(err)
			return
		}
		err = obj.accountRepo.Delete(event.ID)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("[%v] %+v", topic, event)
	default:
		log.Println("no event handler")
	}
}

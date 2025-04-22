package handlers

import (
	"encoding/json"
	cm "go-kafka-bank/consumer/models"
	"go-kafka-bank/events"
	"go-kafka-bank/producer/models"
	"net/http"
	"reflect"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ListAll(mongoClient *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ar := cm.NewAccountRepository(mongoClient)
		results := ar.FindAll()
		c.JSON(http.StatusOK, results)
	}
}

func OpenAccount(eventProducer sarama.SyncProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rBody models.OpenAccount
		if err := c.ShouldBindJSON(&rBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if rBody.AccountHolder == "" || rBody.AccountType <= 0 || rBody.OpeningBalance <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "check your body request (> 0)"})
			return
		}
		event := events.OpenAccountEvent{
			ID: bson.NewObjectID().Hex(),
			AccountHolder:  rBody.AccountHolder,
			AccountType:    rBody.AccountType,
			OpeningBalance: rBody.OpeningBalance,
		}
		if err := Produce(event, eventProducer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
			return
		}
	}
}

func DepositFund(eventProducer sarama.SyncProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rBody models.DepositFund
		if err := c.ShouldBindJSON(&rBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if rBody.ID == "" || rBody.Amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
		event := events.DepositFundEvent{
			ID:     rBody.ID,
			Amount: rBody.Amount,
		}
		if err := Produce(event, eventProducer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
			return
		}
	}
}

func WithdrawFund(eventProducer sarama.SyncProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rBody models.WithdrawFund
		if err := c.ShouldBindJSON(&rBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if rBody.ID == "" || rBody.Amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
		event := events.WithdrawFundEvent{
			ID:     rBody.ID,
			Amount: rBody.Amount,
		}
		if err := Produce(event, eventProducer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
			return
		}
	}
}

func CloseAccount(eventProducer sarama.SyncProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rBody models.WithdrawFund
		if err := c.ShouldBindJSON(&rBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if rBody.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "empty id"})
			return
		}
		event := events.CloseAccountEvent{
			ID: rBody.ID,
		}
		if err := Produce(event, eventProducer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
			return
		}
	}
}

func Produce(event events.Event, eventProducer sarama.SyncProducer) error {
	topic := reflect.TypeOf(event).Name()
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	_, _, err = eventProducer.SendMessage(&msg)
	if err != nil {
		return err
	}
	return nil
}

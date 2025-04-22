package main

import (
	"context"
	"fmt"
	"go-kafka-bank/consumer/db"
	"go-kafka-bank/consumer/models"
	"go-kafka-bank/consumer/services"
	"go-kafka-bank/events"
	"log"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func main() {
	consumer, err := sarama.NewConsumerGroup(viper.GetStringSlice("kafka.servers"), viper.GetString("kafka.group"), nil)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	mongoClient := db.ConnectDatabase()
	defer db.DisconnectDatabase(mongoClient)

	accountRepo := models.NewAccountRepository(mongoClient)
	accountEventHandler := services.NewAccountEventHandler(accountRepo)
	accountConsumerHandler := services.NewConsumerHandler(accountEventHandler)

	fmt.Println("Account consumer started...")
	for {
		consumer.Consume(context.Background(), events.Topics, accountConsumerHandler)
	}
}

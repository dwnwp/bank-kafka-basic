package main

import (
	"context"
	"go-kafka-bank/events"
	"go-kafka-bank/linenoti/services"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
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

	app, err := services.NewApp(
		viper.GetString("line.LINE_CHANNEL_SECRET"),
		viper.GetString("line.LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	port := "8081"
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "healthy",
			"code":    200,
		})
	})
	r.POST("/linebot", app.HandleEvents)

	go func() {
		brokers := viper.GetStringSlice("kafka.servers")
		groupID := viper.GetString("kafka.groupline")
		topics := events.Topics
	
		config := sarama.NewConfig()
		config.Version = sarama.V2_5_0_0
		config.Consumer.Return.Errors = true
		consumer := services.NewConsumer(app)
		
		client, err := sarama.NewConsumerGroup(brokers, groupID, config)
		if err != nil {
			log.Fatalf("Failed to create consumer group: %v", err)
		}
		defer client.Close()
	
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		for {
			if err := client.Consume(ctx, topics, consumer); err != nil {
				log.Fatalf("Error during consumption: %v", err)
			}
			if ctx.Err() != nil {
				break
			}
		}
	}()

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}

}

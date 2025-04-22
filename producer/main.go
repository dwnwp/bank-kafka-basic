package main

import (
	"go-kafka-bank/consumer/db"
	"go-kafka-bank/producer/handlers"
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

	producer, err := sarama.NewSyncProducer(viper.GetStringSlice("kafka.servers"), nil)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	mongoClient := db.ConnectDatabase()
	defer db.DisconnectDatabase(mongoClient)

	port := ":8080"
	r := gin.Default()
	v1 := r.Group("api/v1")
	{
		v1.GET("/healthcheck", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "healthy",
				"code":    200,
			})
		})
		v1.GET("/accounts", handlers.ListAll(mongoClient))
		account := v1.Group("account")
		{
			account.POST("/open", handlers.OpenAccount(producer))
			account.POST("/close", handlers.CloseAccount(producer))
			account.POST("/deposit", handlers.DepositFund(producer))
			account.POST("/withdraw", handlers.WithdrawFund(producer))
		}
	}
	r.Run(port)
}

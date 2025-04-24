package db

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDatabase() *mongo.Client {
	
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connString := viper.GetString("mongodb.CONN_STRING")
	username := viper.GetString("mongodb.MONGO_DB_USERNAME")
	password := viper.GetString("mongodb.MONGO_DB_PASSWORD")

	// MongoDb connection string
	clientOptions := options.Client().ApplyURI(connString)
	// setting auth credentials
	clientOptions.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	log.Println("Connecting to MongoDB...")
	return client
}

func DisconnectDatabase(mongoClient *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mongoClient.Disconnect(ctx); err != nil {
		panic(err)
	}
	log.Println("Disconnected from MongoDB...")
}

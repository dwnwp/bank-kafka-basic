package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BankAccount struct {
	ID            bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AccountHolder string        `json:"accountholder" bson:"accountholder"`
	AccountType   int           `json:"accounttype" bson:"accounttype"`
	Balance       float64       `json:"balance" bson:"balance"`
}

type AccountRepository interface {
	Save(bankAccount BankAccount) error
	Update(bankAccount BankAccount) error
	Delete(id string) error
	FindAll() []BankAccount
	FindByID(id string) (BankAccount, error)
}

type accountRepository struct {
	mongoClient *mongo.Client
}

func NewAccountRepository(mongoClient *mongo.Client) AccountRepository {
	return accountRepository{mongoClient}
}

func (obj accountRepository) Save(bankAccount BankAccount) error {
	collection := obj.mongoClient.Database("UserAccount").Collection("UserAccount")
	_, err := collection.InsertOne(context.TODO(), bankAccount)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (obj accountRepository) Update(bankAccount BankAccount) error {
	collection := obj.mongoClient.Database("UserAccount").Collection("UserAccount")

	filter := bson.M{"_id": bankAccount.ID}

	update := bson.M{
		"$set": bson.M{
			// "accountHolder": bankAccount.AccountHolder,
			// "accountType":   bankAccount.AccountType,
			"balance":       bankAccount.Balance,
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (obj accountRepository) Delete(id string) error {
	resultId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": resultId}
	collection := obj.mongoClient.Database("UserAccount").Collection("UserAccount")
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (obj accountRepository) FindAll() []BankAccount {
	var results []BankAccount
	collection := obj.mongoClient.Database("UserAccount").Collection("UserAccount")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	return results
}

func (obj accountRepository) FindByID(id string) (BankAccount, error) {
	var result BankAccount
	resultId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return BankAccount{}, err
	}
	filter := bson.M{"_id": resultId}
	collection := obj.mongoClient.Database("UserAccount").Collection("UserAccount")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return BankAccount{}, err
	}
	return result, nil
}

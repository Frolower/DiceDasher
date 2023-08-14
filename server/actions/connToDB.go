package actions

import (
	"context"
	"dicedasher/st/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func EstablishConnection() bool {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Создание стека подключений
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
		return false
	}

	db.UserClient = client

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db.UserContext = ctx

	// Подключение к серверу MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return false
	}

	// Проверка подключения
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func CloseConnection(client *mongo.Client, ctx context.Context) {
	client.Disconnect(ctx)
}

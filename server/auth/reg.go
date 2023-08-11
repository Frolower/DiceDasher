package auth

import (
	"context"
	"dicedasher/st/db"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func Reg(c *gin.Context) {
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(400, nil)
		return
	}
	var data map[string]string
	json.Unmarshal(jsonData, &data)

	if data["username"] == "" || data["email"] == "" || data["password_hash"] == "" {
		c.JSON(400, nil)
		return
	}
	// Check password, email, username

	var user = db.User{
		Username:      data["username"],
		Password_hash: data["password_hash"],
		Email:         data["email"],
	}
	//user.GenerateID()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Создание стека подключений
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключение к серверу MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Проверка подключения
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.CreateUser(ctx, client, user)
	if err != nil {
		fmt.Println(err)
		fmt.Println("can't create a user")
	}
	//storage.Users[user.Username] = user

	c.JSON(200, nil)
}

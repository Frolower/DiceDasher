package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Username      string
	Email         string
	Password_hash string
}

func CreateUser(ctx context.Context, client *mongo.Client, user User) error {
	collection := client.Database("diceDasher").Collection("users")

	_, err := collection.InsertOne(ctx, user)

	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}

func ReadUser(ctx context.Context, client *mongo.Client, filter bson.M) (*User, error) {
	collection := client.Database("diceDasher").Collection("users")

	findOptions := options.Find()

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to read users: %v", err)
	}
	defer cur.Close(ctx)

	var users []User

	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("failed to decode user: %v", err)
		}
		users = append(users, user)

		if len(users) > 1 {
			return nil, fmt.Errorf("found more than one user")
		}
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no user found")
	}

	return &users[0], nil
}

func UpdateUser(ctx context.Context, client *mongo.Client, userID primitive.ObjectID, updatedUser User) error {
	collection := client.Database("diceDasher").Collection("users")

	filter := bson.M{"_id": userID}

	update := bson.M{
		"$set": bson.M{
			"username": updatedUser.Username,
			"email":    updatedUser.Email,
			"password": updatedUser.Password_hash,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func DeleteUser(ctx context.Context, client *mongo.Client, userID primitive.ObjectID) error {
	collection := client.Database("diceDasher").Collection("users")

	filter := bson.M{"_id": userID}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

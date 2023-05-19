package storage

import (
	"context"
	"fmt"

	"github.com/Arjune7/booking-go/helper"
	"github.com/Arjune7/booking-go/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


type Message struct {
	Result string
}

func (m *mongoStore) HandleSignUp(name, email, contact, password string) (*types.UserSignUp, error) {
	coll := m.client.Database(m.db.Name()).Collection(m.collection.Name())

	hashPass, err := helper.HandleHashPassword(password)
	if err != nil {
		return nil, err
	}

	doc := bson.M{"name": name, "email": email, "contact": contact, "password": hashPass}
	_, err = coll.InsertOne(context.Background(), doc)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %v", err)
	}
	userCreated := &types.UserSignUp{Name: name, Contact: contact, Password: hashPass, Email: email}
	return userCreated, nil
}

func (m *mongoStore) HandleLogIn(email, password string) (*Message, error) {
	coll := m.client.Database(m.db.Name()).Collection(m.collection.Name())

	filter := bson.M{"email": email}
	user := &types.UserSignUp{}

	err := coll.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &Message{Result: "User not Found"}, nil
		}
		fmt.Println(err)
		return &Message{Result: "Error in finding user"}, err
	}
	result := helper.HandleComparePassword(password, user.Password)
	if !result {
		return &Message{Result: "Password incorrect"}, nil
	}

	token, err := helper.HandleGenerateToken(email, user.Password)
	if err != nil {
		return &Message{Result: "Internal Server error"}, err
	}

	return &Message{Result: token}, nil
}

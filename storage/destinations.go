package storage

import (
	"fmt"
	// "os"

	"github.com/Arjune7/booking-go/helper"
	"github.com/Arjune7/booking-go/types"
	"go.mongodb.org/mongo-driver/bson"

	"context"
	// "github.com/cloudinary/cloudinary-go/v2"
	// "github.com/cloudinary/cloudinary-go/v2/api/admin"
	// "github.com/cloudinary/cloudinary-go/v2/api/uploader"
	// "log"
)

func (m *mongoStore) HandleGetAllDestinations() ([]*types.Destination, error) {
	coll := m.client.Database(m.db.Name()).Collection("Destinations")
	var destinations []*types.Destination

	cursor, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}

	for cursor.Next(context.Background()) {
		destination := &types.Destination{}
		err := cursor.Decode(destination)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %v", err)
		}
		destinations = append(destinations, destination)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}
	return destinations, nil
}

func (m *mongoStore) HandleAddDestination(name, location, price, hostId, rating, placeType, Photos string) (*types.Destination, error) {
	coll := m.client.Database(m.db.Name()).Collection("Destinations")
	placeID := helper.HandleRandomId()
	doc := bson.M{"name": name, "location": location, "price": price, "hostId": hostId, "rating": rating, "placeID": placeID, "placeType": placeType, "Photos": Photos}

	_, err := coll.InsertOne(context.Background(), doc)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document : %v", err)
	}

	destinationAdded := &types.Destination{Name: name, Location: location, Price: price, HostId: hostId, Rating: rating,
		PlaceId: placeID, PlaceType: placeType, Photos: Photos}

	return destinationAdded, nil
}

func (m *mongoStore) HandleGetCategories() ([]*types.Categories, error) {
	coll := m.client.Database(m.db.Name()).Collection("categories-list")

	var categoryList []*types.Categories
	filter := bson.D{}
	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find doc : %v", err)
	}

	for cursor.Next(context.Background()) {
		categoryTypes := &types.Categories{}
		err := cursor.Decode(categoryTypes)
		if err != nil {
			return nil, fmt.Errorf("failed to decode : %v", err)
		}
		categoryList = append(categoryList, categoryTypes)
	}
	return categoryList, nil
}

func (m *mongoStore) HandleAddCategories(name, iconName string) (*types.Categories, error) {
	coll := m.client.Database(m.db.Name()).Collection("categories-list")
	doc := bson.M{"name": name, "iconName": iconName}

	_, err := coll.InsertOne(context.Background(), doc)
	if err != nil {
		return nil, fmt.Errorf("error in uploading file : %v", err)
	}

	return &types.Categories{
		Name:     name,
		IconName: iconName,
	}, nil
}

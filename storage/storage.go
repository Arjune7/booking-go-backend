package storage

import (
	"github.com/Arjune7/booking-go/types"
)

type Storage interface {

	//User Functions

	HandleSignUp(name, email, contact, password string) (*types.UserSignUp, error)
	HandleLogIn(email, password string) (*Message, error)

	//booking and place functions

	HandleGetAllDestinations() ([]*types.Destination, error)
	HandleAddDestination(name, location, price, hostId, rating, placeId, placeType, Photos string) (*types.Destination, error)

	HandleGetCategories() ([]*types.Categories, error)
	HandleAddCategories(name, iconName string) (*types.Categories, error)
}

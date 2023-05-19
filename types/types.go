package types

type Destination struct {
	Name      string `json:"name"`
	Location  string `json:"location"`
	Price     string `json:"price"`
	HostId    string `json:"host_id"`
	Rating    string `json:"rating"`
	PlaceId   string `json:"Place_id"`
	Photos    string `json:"Photos"`
	PlaceType string `json:"Place_type"`
}

type Host struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Contact     string   `json:"contact"`
	HostId      string   `json:"host_id"`
	PlaceId     []string `json:"place_id"`
	Rating      string   `json:"rating"`
	Experience  string   `json:"experience"`
	Reviews     []string `json:"reviews"`
	HostStatus  string   `json:"host_status"`
	Hobbies     []string `json:"hobbies"`
	Description string   `json:"description"`
}

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Contact  string   `json:"contact"`
	Wishlist []string `json:"wishlist"`
	Bookings []string `json:"booking"`
}

type UserSignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Contact  string `json:"contact"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Categories struct {
	Name     string `json:"name"`
	IconName string `json:"iconName"`
}

package server

import (
	"context"
	"encoding/json"
	// "fmt"

	"github.com/Arjune7/booking-go/storage"
	"github.com/Arjune7/booking-go/types"
	"github.com/cloudinary/cloudinary-go/v2"

	//"github.com/cloudinary/cloudinary-go/v2"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/mux"

	// "github.com/joho/godotenv"
	"github.com/rs/cors"
)

type Server struct {
	listenAddr string
	store      storage.Storage
}

type ApiError struct {
	Error string
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content/type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewServer(listenAddr string, store storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) StartServer() error {
	r := mux.NewRouter()

	r.HandleFunc("/start", makeHttpHandleFunc(s.HandleGet)).Methods("GET")
	r.HandleFunc("/signUp", makeHttpHandleFunc(s.HandleUserSignIn)).Methods("POST")
	r.HandleFunc("/login", makeHttpHandleFunc(s.HandleUserLogIn)).Methods("POST")

	//destination Routes
	r.HandleFunc("/destinations/getAll", makeHttpHandleFunc(s.HandleGetAllDestinations)).Methods("GET")
	r.HandleFunc("/destinations/addDestination", makeHttpHandleFunc(s.HandleAddDestination)).Methods("POST")
	r.HandleFunc("/destinations/getCategories", makeHttpHandleFunc(s.HandleGetCategories)).Methods("GET")
	r.HandleFunc("/destinations/addCategories", makeHttpHandleFunc(s.HandleAddCategories)).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Change this to your allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	})

	// Wrap your router with the CORS middleware
	handler := c.Handler(r)

	return http.ListenAndServe(s.listenAddr, handler)
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, "HELLO I AM HERE")
}

func (s *Server) HandleUserSignIn(w http.ResponseWriter, r *http.Request) error {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	newUser := &types.UserSignUp{}
	err = json.Unmarshal(body, newUser)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}

	if newUser.Name == "" || newUser.Contact == "" || newUser.Password == "" || newUser.Email == "" {
		return writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	user, err := s.store.HandleSignUp(newUser.Name, newUser.Email, newUser.Contact, newUser.Password)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}

	return writeJSON(w, http.StatusOK, user)
}

func (s *Server) HandleUserLogIn(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}
	user := &types.UserLogin{}
	err = json.Unmarshal(body, user)

	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}
	res, err := s.store.HandleLogIn(user.Email, user.Password)
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	return writeJSON(w, http.StatusOK, res)
}

//Destination functions

func (s *Server) HandleGetAllDestinations(w http.ResponseWriter, r *http.Request) error {
	destinations, err := s.store.HandleGetAllDestinations()
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}

	return writeJSON(w, http.StatusOK, destinations)
}

func (s *Server) HandleAddDestination(w http.ResponseWriter, r *http.Request) error {

	err := r.ParseMultipartForm(32 << 20) // Parse up to 32MB of data
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}

	// Access individual form values by name
	name := r.FormValue("name")
	location := r.FormValue("location")
	price := r.FormValue("price")
	hostId := r.FormValue("hostId")
	rating := r.FormValue("rating")
	placeType := r.FormValue("placeType")

	// err = godotenv.Load()
	// if err != nil {
	// 	// Handle the error if the .env file couldn't be loaded
	// 	fmt.Printf("Error loading .env file: %v\n", err)
	// 	return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	// }

	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUD_NAME"), os.Getenv("API_KEY"), os.Getenv("API_SECRET"))
	if err != nil {

		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	// Create a pipe to stream the file data
	pr, pw := io.Pipe()

	// Create a goroutine to stream the file data to the pipe writer
	go func() {
		defer func(pw *io.PipeWriter) {
			err := pw.Close()
			if err != nil {
				return
			}
		}(pw)

		// Copy the file data to the pipe writer
		_, err := io.Copy(pw, file)
		if err != nil {
			err := writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
			if err != nil {
				return
			}
		}
	}()

	// Prepare the upload parameters
	uploadParams := uploader.UploadParams{
		PublicID: handler.Filename,
	}

	// Perform the upload to Cloudinary
	result, err := cld.Upload.Upload(context.Background(), pr, uploadParams)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}
	destination, err := s.store.HandleAddDestination(name, location, price, hostId, rating, placeType, result.SecureURL)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}

	return writeJSON(w, http.StatusOK, destination)
}

func (s *Server) HandleGetCategories(w http.ResponseWriter, r *http.Request) error {
	category, err := s.store.HandleGetCategories()
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}
	return writeJSON(w, http.StatusOK, category)
}

func (s *Server) HandleAddCategories(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	newCategory := &types.Categories{}
	err = json.Unmarshal(body, newCategory)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}

	result, err := s.store.HandleAddCategories(newCategory.Name, newCategory.IconName)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}

	return writeJSON(w, http.StatusOK, result)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Arjune7/booking-go/server"
	"github.com/Arjune7/booking-go/storage"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	listener := flag.String("listenAddr", "0.0.0.0:"+port, "newServer running at port 8000")
	flag.Parse()
	ctx := context.Background()

	store, err := storage.NewDatabase(ctx, "mongodb+srv://"+os.Getenv("USERNAME")+":"+os.Getenv("PASSWORD")+"@cluster0.imrset9.mongodb.net/?retryWrites=true&w=majority", "booking-go", "user-info")
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	defer func() {
		err := store.CloseStore(ctx)
		if err != nil {
			fmt.Printf("Error closing store: %v\n", err)
			return
		}
	}()

	newServer := server.NewServer(*listener, store)
	fmt.Println("newServer running at port", *listener)

	log.Fatal(newServer.StartServer())
}

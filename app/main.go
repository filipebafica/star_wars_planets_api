package main

import (
	"fmt"
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/filipebafica/star_wars_planets_api/data"
	"github.com/filipebafica/star_wars_planets_api/handlers"
)

func main() {
	fmt.Println("application has started...")

	// define a context that carries the time that will be used as limit to db connection attempt
	// define a function callback in case timeout is reached
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// will release the resources if application hits time out
	// this is going to be executed when the function reaches the end of its scope
	defer cancel()
	// handle db connection
	data.ConnectDB(ctx)
	defer data.DisconnectDB(ctx)

	// define a router to handle requests
	router := mux.NewRouter()
	// define the specific handlers to each request and endpoint
	router.HandleFunc("/planet", handlers.CreatePlanetEndPoint).Methods("POST")
	// loop to listen and respond for requests
	http.ListenAndServe(":8000", router)
}

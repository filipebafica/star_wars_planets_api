package handlers

import (
	"context"
	"net/http"
	"time"
	"encoding/json"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	"github.com/filipebafica/star_wars_planets_api/data"
)

func CreatePlanetEndPoint(response http.ResponseWriter, request *http.Request) {
	// define the response content-type as json
	response.Header().Set("content-type", "application/json")

	// define a variable that will receive the request data
	var planet data.Planet

	// decode request data into planet variable
	json.NewDecoder(request.Body).Decode(&planet)

	// define a context that carries the time that will be used as limit to db operation attempt
	// skip the callback function since errors will be handled if find does not match
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// try to insert data into db
	// if fails, messege error is returned
	result, err := data.Collection.InsertOne(ctx, planet)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	// encode the insertion message return into json formart and append it to the response
	json.NewEncoder(response).Encode(result)
}

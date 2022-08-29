package handlers

import (
	"context"
	"net/http"
	"time"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
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

func GetPlanetsEndPoint(response http.ResponseWriter, request *http.Request) { 
	// define the response content-type as json
	response.Header().Set("content-type", "application/json")

	// define a 'dynamically-sized array' that will receive queries from db
	var planets []data.Planet

	// define a context that carries the time that will be used as limit to db operation attempt
	// skip the callback function since errors will be handled if find does not match
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// get a cursor 'pointer' to the entries in db with an empty filter 'bson.M{}'
	// if fails, messege error is returned
	cursor, err := data.Collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	//if not fail it will close cursor at the end of fucntion's scope
	defer cursor.Close(ctx)

	// iteraates through cursor and appen to the people slice
	// if fails, messege error is returned
	if err := cursor.All(ctx, &planets); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	// encode people slice into json format and append it to the response
	json.NewEncoder(response).Encode(planets)
}

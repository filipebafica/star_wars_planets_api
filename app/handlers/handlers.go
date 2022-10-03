package handlers

import (
	"context"
	"net/http"
	"time"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/filipebafica/star_wars_planets_api/data"
)

func CreatePlanetEndPoint(response http.ResponseWriter, request *http.Request) {
	// define the response content-type as json
	response.Header().Set("content-type", "application/json")

	// define a variable that will receive the request data
	var planet data.Planet

	// decode request data into planet variable
	json.NewDecoder(request.Body).Decode(&planet)
	defer request.Body.Close()

	// validate request
	if err := planet.Validate(); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	// define a context that carries the time that will be used as limit to db operation attempt
	// skip the callback function since errors will be handled if find does not match
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// add how many films the planet was in
	swapiResults := data.GetFilmsPlanetWasIn(planet.Name)
	if swapiResults.Count == 0 || swapiResults.Results[0].Name != planet.Name {
		planet.Films = 0
	} else {
		planet.Films = len(swapiResults.Results[0].Films)
	}

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

	// define a 'dynamically-sized array' that will receive queries from db
	var planets []data.Planet

	// iteraates through cursor and append to the planets slice
	// if fails, messege error is returned
	if err := cursor.All(ctx, &planets); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	// encode planets slice into json format and append it to the response
	json.NewEncoder(response).Encode(planets)
}

func GetPlanetEndPoint(response http.ResponseWriter, request *http.Request) {
	// define the response content-type as json
	response.Header().Set("content-type", "application/json")

	// get and check params from query string
	query := request.URL.Query()
		if query.Get("id") == "" && query.Get("nome") == "" {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte(`{"message": "` + `Resource Not Found` + `"}`))
			return
		}

	var filter bson.M
	if query.Get("id") != "" {
		v, _ := primitive.ObjectIDFromHex(query.Get("id"))
		filter = bson.M{"_id":v}
	} else {
		v := query.Get("nome")
		filter = bson.M{"name":v}
	}

	// define a context that carries the time that will be used as limit to db operation attempt
	// skip the callback function since errors will be handled if find does not match
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// get a cursor 'pointer' to the entries in db with filter
	// if fails, messege error is returned
	cursor, err := data.Collection.Find(ctx, filter)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	// if value not found
	if cursor.RemainingBatchLength() == 0 {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "` + `Resource Not Found` + `"}`))
		return
	}
	//if not fail it will close cursor at the end of fucntion's scope
	defer cursor.Close(ctx)

	// define a 'dynamically-sized array' that will receive queries from db
	var planet []data.Planet

	// iteraates through cursor and appen to the planet slice
	// if fails, messege error is returned
	if err := cursor.All(ctx, &planet); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	// encode planet slice into json format and append it to the response
	json.NewEncoder(response).Encode(planet)
}

func DeletePlanetEndPoint(response http.ResponseWriter, request *http.Request) { 
	// define the response content-type as json
	response.Header().Set("content-type", "application/json")

	// get and check params from query string
	query := request.URL.Query()
		if query.Get("id") == "" {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte(`{"message": "` + `Resource Not Found` + `"}`))
			return
		}

	v, _ := primitive.ObjectIDFromHex(query.Get("id"))
	filter := bson.M{"_id":v}

	// define a context that carries the time that will be used as limit to db operation attempt
	// skip the callback function since errors will be handled if find does not match
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)


	// try to find a planet by its ID and delete it
	var deletedDocument data.Planet
	err := data.Collection.FindOneAndDelete(ctx, filter).Decode(&deletedDocument)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "` + `Resource Not Found` + `"}`))
		return
	} else {
		response.WriteHeader(http.StatusOK)
		response.Write([]byte(`{"message": "` + `Resource Has Been Deleted` + `"}`))
		return
	}
}

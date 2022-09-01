package data

import (
	"context"
	"net/http"
	"log"
	"encoding/json"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/go-playground/validator/v10"
)

var Client *mongo.Client
var Collection *mongo.Collection

type Planet struct {
	ID						primitive.ObjectID	`json:"_id,omitempty" bson:"_id,omitempty"`
	Name					string				`json:"nome,omitempty" bson:"name,omitempty" validate:"required"`
	Climate					string				`json:"clima,omitempty" bson:"climate,omitempty" validate:"required"`
	Terrain					string				`json:"terreno,omitempty" bson:"terrain,omitempty" validate:"required"`
	Films					int					`json:"filmes" bson:"films,omitempty"`
}

type SWAPIResults struct {
	Count		int			`json:"count"`
	Results		[]struct {
		Name	string		`json:"name"`
		Films	[]string	`json:"films"`
	}	`json:"results"`
}

func GetFilmsPlanetWasIn(planet_name string) SWAPIResults {
	resp, err := http.Get("https://swapi.dev/api/planets/?search=" + planet_name)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var swapiResults SWAPIResults
	json.NewDecoder(resp.Body).Decode(&swapiResults)

	return swapiResults
}

func ConnectDB(ctx context.Context) {
	// try to establish connection with the mongodb
	// if fails, program will finish execution
	Client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://user:password@mongodb:27017"))
	if err != nil {
		panic(err)
	}

	// define which database and which collection 'table' is going to be accessed
	Collection = Client.Database("starwarsdb").Collection("planets")
}

func DisconnectDB(ctx context.Context) {
	// will disconnect from database at the end of program execution
	if err := Client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (planet *Planet) Validate() error {
	validate := validator.New()
	return validate.Struct(planet)
}

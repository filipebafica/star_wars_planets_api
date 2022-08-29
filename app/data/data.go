package data

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Collection *mongo.Collection

type Planet struct {
	ID						primitive.ObjectID	`json:"_id,omitempty" bson:"_id,omitempty"`
	Nome					string				`json:"nome,omitempty" bson:"nome,omitempty"`
	Clima					string				`json:"clima,omitempty" bson:"clima,omitempty"`
	Terreno					string				`json:"terreno,omitempty" bson:"terreno,omitempty"`
	Aparicoes_Em_Filmes		int64				`json:"aparicoes_em_filmes,omitempty" bson:"aparicoes_em_filmes,omitempty"`
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

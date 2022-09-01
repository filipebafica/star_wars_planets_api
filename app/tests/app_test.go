package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"context"
	"time"

	"github.com/gorilla/mux"
	"github.com/steinfletcher/apitest"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/filipebafica/star_wars_planets_api/data"
	"github.com/filipebafica/star_wars_planets_api/handlers"

)

func TestInsertPlanet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data.ConnectDB(ctx, "unitTests")
	defer data.DisconnectDB(ctx)

	r := mux.NewRouter()
	r.HandleFunc("/v1/planeta", handlers.CreatePlanetEndPoint).Methods("POST")
	ts := httptest.NewServer(r)
	defer ts.Close()
	t.Run("Valid Insertion", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/v1/planeta").
			JSON(`{"nome": "Tatooine", "clima": "arido", "terreno": "deserto"}`).
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	t.Run("Invalid Insertion", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/v1/planeta").
			JSON(`{"clima": "arido", "terreno": "deserto"}`).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestSearchPlanet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data.ConnectDB(ctx, "unitTests")
	defer data.DisconnectDB(ctx)

	r := mux.NewRouter()
	r.HandleFunc("/v1/planeta", handlers.CreatePlanetEndPoint).Methods("POST")
	r.HandleFunc("/v1/planeta", handlers.GetPlanetEndPoint).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()
	t.Run("Valid Insertion", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/v1/planeta").
			JSON(`{"nome": "Tatooine", "clima": "arido", "terreno": "deserto"}`).
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	t.Run("Valid Search By Name", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/v1/planeta?nome=Tatooine").
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	t.Run("Invalid Search By Name", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/v1/planeta?nome=test").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	var result data.Planet
	err := data.Collection.FindOne(ctx, bson.D{{"nome", "Tatooine"}}).Decode(&result)
	if err != nil {
		return
	}

	t.Run("Valid Search By ID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/v1/planeta?id=" + string(result.ID.Hex())).
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	t.Run("Invalid Search By ID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/v1/planeta?id=test").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})
}

func TestDeletePlanet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data.ConnectDB(ctx, "unitTests")
	defer data.DisconnectDB(ctx)

	r := mux.NewRouter()
	r.HandleFunc("/v1/planeta", handlers.CreatePlanetEndPoint).Methods("POST")
	r.HandleFunc("/v1/planeta", handlers.DeletePlanetEndPoint).Methods("DELETE")
	ts := httptest.NewServer(r)
	defer ts.Close()
	t.Run("Valid Insertion", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/v1/planeta").
			JSON(`{"nome": "Tatooine", "clima": "arido", "terreno": "deserto"}`).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	var result data.Planet
	err := data.Collection.FindOne(ctx, bson.D{{"nome", "Tatooine"}}).Decode(&result)
	if err != nil {
		return
	}

	t.Run("Valid Deletion", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Delete("/v1/planeta?id=" + string(result.ID.Hex())).
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	t.Run("Invalid Deletion", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Delete("/v1/planeta?id=test").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})
}


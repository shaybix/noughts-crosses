package main

import (
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	storage := &BoltDB{
		db: db,
	}

	game := &Game{
		db: storage,
	}

	router.HandleFunc("/game", game.CreateGameHandler).Methods("POST")
	router.HandleFunc("/game/{id}/join", game.JoinPlayerHandler).Methods("POST")
	router.HandleFunc("/game/{id}", game.GetGameHandler).Methods("GET")
	router.HandleFunc("/game/{id}/nought", game.SetNoughtHandler).Methods("POST")
	router.HandleFunc("/game/{id}/cross", game.SetCrossHandler).Methods("POST")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 20,
	}

	log.Fatal(server.ListenAndServe())
}

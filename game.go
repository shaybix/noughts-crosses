package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Game represents the entire game between two players.
type Game struct {
	ID           int      `json:"id"`
	FirstPlayer  *User    `json:"first_player"`
	SecondPlayer *User    `json:"second_player"`
	Noughts      []Nought `json:"noughts"`
	Crosses      []Cross  `json:"crosses"`
	Finished     bool     `json:"finished"`
	Winner       *User    `json:"winner"`
	db           Database
}

// User is an individual player playing the game.
type User struct {
	ID       int  `json:"id"`
	IsNought bool `json:"is_nought"`
}

// Nought represents an individual nought in game with an x and y coordinates on the board.
type Nought struct {
	X      int   `json:"x"`
	Y      int   `json:"y"`
	Player *User `json:"player"`
}

// Cross represents an individual cross in game with an x and y coordinates on the board.
type Cross struct {
	X      int   `json:"x"`
	Y      int   `json:"y"`
	Player *User `json:"player"`
}

// CreateGameHandler creates a new game and returns that initialised game in json format.
func (g *Game) CreateGameHandler(w http.ResponseWriter, r *http.Request) {

	user := &User{
		ID:       rand.Int(),
		IsNought: true,
	}
	g.FirstPlayer = user

	if err := g.db.createGame(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// JoinPlayerHandler registers a second player to join a game that the first player created.
func (g *Game) JoinPlayerHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	gameID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = g.db.getGame(gameID, g); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if g.SecondPlayer != nil {
		http.Error(w, fmt.Sprint("game already has two players"), http.StatusBadRequest)
		return
	}

	user := &User{
		ID: rand.Int(),
	}

	g.SecondPlayer = user

	if err = g.db.saveGame(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// GetGameHandler gets an existing game with the ID given in /game/{id} in the route and returns the status of such a game.
func (g *Game) GetGameHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	gameID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = g.db.getGame(gameID, g); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// GetWinnerHandler checks if there is a winner and if there is a winner returns
func (g *Game) GetWinnerHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	if !g.Finished {
		http.Error(w, fmt.Sprint("There is no winner"), http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(g.Winner); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// SetNoughtHandler sets a nought in a single move made by a player.
func (g *Game) SetNoughtHandler(w http.ResponseWriter, r *http.Request) {

	var n Nought

	vars := mux.Vars(r)
	gameID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = g.db.getGame(gameID, g); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if g.Finished == true {
		http.Error(w, fmt.Sprint("game is already finished"), http.StatusBadRequest)
		return

	}

	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if n.Player == nil {
		http.Error(w, fmt.Sprint("the player is not set for nought"), http.StatusBadRequest)
		return
	}

	g.Noughts = append(g.Noughts, n)

	if err = g.db.saveGame(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(g.Noughts) >= 3 {
		if hasNoughtsWon(g.Noughts) {
			g.Finished = true
			g.Winner = n.Player
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func hasNoughtsWon(n []Nought) bool {

	//TODO: need to traverse through the slice of noughts and find a way to check
	// if there is a consecutive 3 noughts and if so then return true

	return false
}

// SetCrossHandler sets a cross in a single move made by a player.
func (g *Game) SetCrossHandler(w http.ResponseWriter, r *http.Request) {

	var c Cross

	vars := mux.Vars(r)
	gameID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = g.db.getGame(gameID, g); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if g.Finished == true {
		http.Error(w, fmt.Sprint("game is already finished"), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if c.Player == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.Crosses = append(g.Crosses, c)

	if err = g.db.saveGame(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(g.Crosses) >= 3 {
		if hasCrossesWon(g.Crosses) {
			g.Finished = true
			g.Winner = c.Player
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func hasCrossesWon(c []Cross) bool {

	//TODO: need to traverse through the slice of crosses and find a way to check
	// if there is a consecutive 3 crosses and if so then return true

	return false
}

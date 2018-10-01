package main

import (
	"encoding/json"
	"math/rand"
	"strconv"

	"github.com/boltdb/bolt"
)

// Database is an interface
type Database interface {
	createGame(g *Game) error
	saveGame(g *Game) error
	getGame(id int, g *Game) error
}

// BoltDB is a database type that satisfies the Database interface so in the future one can easily
// swap it for another choice of database, so long it satisfies the Database interface.
type BoltDB struct {
	db *bolt.DB
}

func (blt BoltDB) createGame(g *Game) error {

	err := blt.db.Update(func(tx *bolt.Tx) error {
		g.ID = rand.Int()
		b, err := tx.CreateBucketIfNotExists([]byte("games"))
		if err != nil {
			return err
		}
		buf, err := json.Marshal(g)
		if err != nil {
			return err
		}
		b.Put([]byte(strconv.Itoa(g.ID)), buf)
		return nil
	})

	return err
}

func (blt BoltDB) saveGame(g *Game) error {
	var err error

	err = blt.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("games"))
		buf, err := json.Marshal(g)
		if err != nil {
			return err
		}
		err = b.Put([]byte(strconv.Itoa(g.ID)), buf)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (blt BoltDB) getGame(id int, g *Game) error {
	var err error
	err = blt.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("games"))
		buf := b.Get([]byte(strconv.Itoa(id)))
		err := json.Unmarshal(buf, g)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

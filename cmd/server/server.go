package main

import (
	"log"

	"github.com/chinnaxs/go_beer/internal/pkg/db"
)

func main() {
	db, err := db.NewFileDB("test.json")
	if err != nil {
		log.Println(err)
		return
	}
	b, err := db.Beer("test2")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(b)
	err = db.RemoveBeer("test2")
	if err != nil {
		log.Println(err)
		return
	}
}

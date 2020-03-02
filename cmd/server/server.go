package main

import (
	"log"
	"net"
	"net/http"

	"github.com/chinnaxs/go_beer/internal/pkg/db"
	"github.com/chinnaxs/go_beer/internal/pkg/server"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("was not able to parse tcp Addr")
	}
	fileDB, err := db.NewFileDB("beerstore.json")
	if err != nil {
		log.Fatalf("was not able to connect to database: %s", err)
	}
	s := server.NewServer(tcpAddr, fileDB)
	s.Start()

	log.Fatal(http.ListenAndServe(":8080", nil))
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

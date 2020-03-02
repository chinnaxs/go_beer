package server

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/chinnaxs/go_beer/internal/pkg/beverage"
	"github.com/chinnaxs/go_beer/internal/pkg/db"
)

type Server struct {
	addr *net.TCPAddr
	db   db.DB
}

func NewServer(addr *net.TCPAddr, db db.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Start() {
	http.HandleFunc("/beers", s.beersHandler)
	http.HandleFunc("/beers/", s.beerHandler)
	log.Printf("Start listening on %s", s.addr.String())
	log.Fatal(http.ListenAndServe(s.addr.String(), nil))
}

func (s *Server) beersHandler(w http.ResponseWriter, _ *http.Request) {
	beers, err := s.db.Beers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	enc, err := json.Marshal(beers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(enc)
	if err != nil {
		log.Println("Was not able to response")
	}
}

func (s *Server) beerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetBeer(s, w, r)
	case http.MethodPut:
		handlePutBeer(s, w, r)
	case http.MethodDelete:
		handleDeleteBeer(s, w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleGetBeer(s *Server, w http.ResponseWriter, r *http.Request) {
	beerName, err := extractBeerName(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	beer, err := s.db.Beer(beerName)
	if err != nil {
		if strings.HasPrefix(err.Error(), "no beer named") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	enc, err := json.Marshal(beer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(enc)
	if err != nil {
		log.Println("Was not able to response")
	}
}

func handlePutBeer(s *Server, w http.ResponseWriter, r *http.Request) {
	beerName, err := extractBeerName(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var beer beverage.Beer
	err = json.NewDecoder(r.Body).Decode(&beer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if beer.Name != beerName {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = s.db.Beer(beerName)
	if err != nil {
		if strings.HasPrefix(err.Error(), "no beer named") {
			createBeer(s, w, beer)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	updateBeer(s, w, beer)
}

func createBeer(s *Server, w http.ResponseWriter, beer beverage.Beer) {
	err := s.db.AddBeer(beer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//TODO CFE set location
	w.WriteHeader(http.StatusCreated)
}

func updateBeer(s *Server, w http.ResponseWriter, beer beverage.Beer) {
	err := s.db.UpdateBeer(beer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteBeer(s *Server, w http.ResponseWriter, r *http.Request) {
	beerName, err := extractBeerName(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.db.RemoveBeer(beerName)
	if err != nil {
		if strings.HasPrefix(err.Error(), "no beer named") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func extractBeerName(url string) (string, error) {
	//TODO validate input
	return strings.Split(url, "/")[2], nil
}

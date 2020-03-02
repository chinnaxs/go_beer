package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/chinnaxs/go_beer/internal/pkg/beverage"
)

type FileDB struct {
	beers    map[string]*beverage.Beer
	fileName string
	mux      sync.RWMutex
}

func NewFileDB(fileName string) (DB, error) {
	db := &FileDB{
		beers:    make(map[string]*beverage.Beer),
		fileName: fileName,
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &db.beers); err != nil {
		return nil, err
	}
	return db, nil
}

func (m *FileDB) Beers() ([]beverage.Beer, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	beers := []beverage.Beer{}
	for _, v := range m.beers {
		beers = append(beers, *v)
	}
	return beers, nil
}

func (m *FileDB) Beer(name string) (beverage.Beer, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	beer, ok := m.beers[name]
	if !ok {
		return beverage.Beer{}, fmt.Errorf("no beer named %s exists", name)
	}
	return *beer, nil
}

func (m *FileDB) AddBeer(beer beverage.Beer) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.beers[beer.Name]; ok {
		return fmt.Errorf("Beer with name %s is already in the store", beer.Name)
	}
	m.beers[beer.Name] = &beer
	enc, err := json.Marshal(m.beers)
	if err != nil {
		delete(m.beers, beer.Name)
		return err
	}
	if err := ioutil.WriteFile(m.fileName, enc, 0644); err != nil {
		delete(m.beers, beer.Name)
		return err
	}
	return nil
}

func (m *FileDB) UpdateBeer(beer beverage.Beer) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	oldBeer, ok := m.beers[beer.Name]
	if !ok {
		return fmt.Errorf("Beer with name %s is not in the store", beer.Name)
	}
	m.beers[beer.Name] = &beer
	enc, err := json.Marshal(m.beers)
	if err != nil {
		m.beers[beer.Name] = oldBeer
		return err
	}
	if err := ioutil.WriteFile(m.fileName, enc, 0644); err != nil {
		m.beers[beer.Name] = oldBeer
		return err
	}
	return nil
}

func (m *FileDB) RemoveBeer(name string) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	oldBeer, ok := m.beers[name]
	if !ok {
		return fmt.Errorf("Beer with name %s is not in the store", name)
	}
	delete(m.beers, name)
	enc, err := json.Marshal(m.beers)
	if err != nil {
		m.beers[oldBeer.Name] = oldBeer
		return err
	}
	if err := ioutil.WriteFile(m.fileName, enc, 0644); err != nil {
		m.beers[oldBeer.Name] = oldBeer
		return err
	}
	return nil
}

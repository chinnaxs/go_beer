package db

import (
	"fmt"
	"sync"

	"github.com/chinnaxs/go_beer/internal/pkg/beverage"
)

type MemoryDB struct {
	beers map[string]*beverage.Beer
	mux   sync.RWMutex
}

func NewMemoryDB() DB {
	return &MemoryDB{beers: make(map[string]*beverage.Beer)}
}

func (m *MemoryDB) Beers() ([]beverage.Beer, error) {
	m.mux.RLocker()
	defer m.mux.RUnlock()
	beers := []beverage.Beer{}
	for _, v := range m.beers {
		beers = append(beers, *v)
	}
	return beers, nil
}

func (m *MemoryDB) Beer(name string) (beverage.Beer, error) {
	m.mux.RLocker()
	defer m.mux.RUnlock()
	beer, ok := m.beers[name]
	if !ok {
		return beverage.Beer{}, fmt.Errorf("no beer named %s exists", name)
	}
	return *beer, nil
}

func (m *MemoryDB) AddBeer(beer beverage.Beer) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.beers[beer.Name]; ok {
		return fmt.Errorf("Beer with name %s is already in the store", beer.Name)
	}
	m.beers[beer.Name] = &beer
	return nil
}

func (m *MemoryDB) UpdateBeer(beer beverage.Beer) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.beers[beer.Name]; !ok {
		return fmt.Errorf("Beer with name %s is not in the store", beer.Name)
	}
	m.beers[beer.Name] = &beer
	return nil
}

func (m *MemoryDB) RemoveBeer(name string) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.beers[name]; !ok {
		return fmt.Errorf("Beer with name %s is not in the store", name)
	}
	delete(m.beers, name)
	return nil
}

package db

import "github.com/chinnaxs/go_beer/internal/pkg/beverage"

type DB interface {
	//Beers returns a list of available beers in the store. Returns an error if the beer list cannot
	//be accessed.
	Beers() ([]beverage.Beer, error)
	//Beer returns information about this beer. Returns an error if the beer list cannot
	//be accessed.
	Beer(name string) (beverage.Beer, error)
	//AddBeer adds beer to the store. Returns an error if a beer with the same name is already in
	//the store or if the information about the beer cannot be persisted.
	AddBeer(beer beverage.Beer) error
	//UpdateBeer changes information about beer. Returns an error if beer is not part of the store
	//or if the new information about the beer cannot be persisted.
	UpdateBeer(beer beverage.Beer) error
	//RemoveBeer removes this beer from the store. Returns an error if the beer is not part of the
	//store or if the beer list cannot be accessed.
	RemoveBeer(name string) error
}

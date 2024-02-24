package handlers

import (
	"sync"
)

type product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var (
	UID int

	Products = make(map[int]product)

	globalMuUID sync.Mutex
)

func addProductToDB(p *newProduct) product {
	globalMuUID.Lock()
	defer globalMuUID.Unlock()
	newProductForDB := product{Id: UID, Name: p.Name, Description: p.Description}
	Products[UID] = newProductForDB
	UID++

	return newProductForDB
}

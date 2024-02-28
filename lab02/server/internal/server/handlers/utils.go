package handlers

import (
	"sync"
)

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

func updateProductInDB(p newProduct, productFromList product) product {
	globalMuUID.Lock()
	defer globalMuUID.Unlock()

	var updProduct product
	updProduct.Id = productFromList.Id
	if p.Name != "" {
		updProduct.Name = p.Name
	} else {
		updProduct.Name = productFromList.Name
	}

	if p.Description != "" {
		updProduct.Description = p.Description
	} else {
		updProduct.Description = productFromList.Description
	}

	return updProduct
}

func deleteProductFromDB(idOfProduct int) {
	globalMuUID.Lock()
	defer globalMuUID.Unlock()

	delete(Products, idOfProduct)
}

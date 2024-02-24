package server

import (
	"log"
	"net/http"

	handlers "lab_2/internal/server/handlers"
)

func RunServer() {
	// init handlers
	http.HandleFunc("/product", handlers.AddProduct)
	http.HandleFunc("/product/{product_id}", handlers.ProductsOperations)
	http.HandleFunc("/products/", handlers.GetAllProducts)

	// run server
	log.Printf("<--------SERVER RUNNING---------->")
	http.ListenAndServe(":8080", nil)
}

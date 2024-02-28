package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Write([]byte(fmt.Sprintf("method must be get!")))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var allProducts []product
	for _, p := range Products {
		allProducts = append(allProducts, p)
	}

	productsJSON, err := json.Marshal(allProducts)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("can not marsahl new item:%v", err)))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Write([]byte(fmt.Sprintf("%v", string(productsJSON))))
	w.WriteHeader(http.StatusOK)
}

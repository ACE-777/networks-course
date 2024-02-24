package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type newProduct struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("method must be POST")))

		return
	}

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	var p newProduct
	if err := decode.Decode(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("can not decode json message, invalid format: %v", err)))

		return
	}

	newProductForResponse := addProductToDB(&p)
	newProductJSON, err := json.Marshal(newProductForResponse)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("can not marsahl new item:%v", err)))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Write([]byte(fmt.Sprintf("%v", string(newProductJSON))))
	w.WriteHeader(http.StatusOK)
}

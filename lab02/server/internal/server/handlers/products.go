package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type newProduct struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ProductsOperations(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
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

		return

	case http.MethodGet:
		uri := strings.Split(r.URL.Path, "/")
		idOfProduct, err := strconv.Atoi(uri[2])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error while getting integer from string: %v", err)))

			return
		}

		productFromList, ok := Products[idOfProduct]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("undefined product in db")))

			return
		}

		productFromListJSON, err := json.Marshal(productFromList)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("can not marsahl item from db: %v", err)))
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(fmt.Sprintf("%v", string(productFromListJSON))))
		w.WriteHeader(http.StatusOK)

		return

	case http.MethodPut:
		uri := strings.Split(r.URL.Path, "/")
		idOfProduct, err := strconv.Atoi(uri[2])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error while getting integer from string: %v", err)))

			return
		}

		productFromList, ok := Products[idOfProduct]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("undefined product in db")))

			return
		}

		decode := json.NewDecoder(r.Body)
		decode.DisallowUnknownFields()

		var p newProduct
		if err = decode.Decode(&p); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("can not decode json message, invalid format: %v", err)))

			return
		}

		updProduct := updateProductInDB(p, productFromList)

		Products[idOfProduct] = updProduct

		productFromListJSON, err := json.Marshal(updProduct)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("can not marsahl item from db: %v", err)))
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(fmt.Sprintf("%v", string(productFromListJSON))))
		w.WriteHeader(http.StatusOK)

		return

	case http.MethodDelete:
		uri := strings.Split(r.URL.Path, "/")
		idOfProduct, err := strconv.Atoi(uri[2])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error while getting integer from string: %v", err)))

			return
		}

		productFromList, ok := Products[idOfProduct]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("undefined product in db")))

			return
		}

		deleteProductFromDB(idOfProduct)

		productFromListJSON, err := json.Marshal(productFromList)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("can not marsahl item from db: %v", err)))
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(fmt.Sprintf("%v", string(productFromListJSON))))
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("invalid method")))

		return

	}
}

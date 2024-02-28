package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
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
	Icon        string `json:"icon"`
}

func ProductsOperations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		uri := strings.Split(r.URL.Path, "/")
		if len(uri) > 3 {
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

			icon, _, err := r.FormFile("icon")
			defer icon.Close()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("can't receive icon")))

			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("can not load icon")))

				return
			}

			localIcon, err := os.Create(path.Join("internal", "server", "handlers", "icons", strconv.Itoa(idOfProduct)+".png"))
			defer localIcon.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("can not open icon for writting")))

				return
			}

			_, err = io.Copy(localIcon, icon)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("can not copy icon for local file system")))

				return
			}

			productFromList.Icon = path.Join("internal", "server", "handlers", "icons",
				strconv.Itoa(idOfProduct)+".png")
			Products[idOfProduct] = productFromList

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("successfully upload icon"))

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

		if len(uri) > 4 {
			w.Write([]byte(productFromList.Icon))
			w.WriteHeader(http.StatusOK)

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

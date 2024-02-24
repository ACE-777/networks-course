package handlers

import "net/http"

func ProductsOperations(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodPut || r.Method == http.MethodDelete {

		return
	}

	w.WriteHeader(http.StatusInternalServerError)

	return
}

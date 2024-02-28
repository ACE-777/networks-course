package handlers

import "net/http"

func ProductIcons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		return
	}

	if r.Method != http.MethodGet {

		return
	}

}

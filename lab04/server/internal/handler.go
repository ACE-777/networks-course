package internal

import (
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	https       = "https://"
	contentType = "application/json"
)

var journal = make(map[string]int)

type ProxyServer struct{}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s\n", r.Method, r.URL.String())

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		return
	}

	destinationURL := *r.URL
	destinationURL.Scheme = "http"
	destinationURL.Host = r.Host
	proxyRequest := &http.Request{
		Method: r.Method,
		URL:    &destinationURL,
		Header: r.Header,
		Body:   r.Body,
	}

	url := strings.TrimPrefix(proxyRequest.URL.Path, "/")
	if r.Method == http.MethodGet {
		resp, err := http.Get(https + url)
		if err != nil {
			http.Error(w, "Cant do GET", http.StatusInternalServerError)
			log.Printf("can done get request, err: %v", err)

			return
		}

		log.Printf("Response status code: %v\n", resp.Status)

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)

		journal[url] = resp.StatusCode

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("Error copying response body: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}

		log.Printf("success execute handler")

		return
	}

	resp, err := http.Post(https+url, contentType, r.Body)
	if err != nil {
		http.Error(w, "Cant do POST", http.StatusInternalServerError)
		log.Printf("can done get request, err: %v", err)

		return
	}

	log.Printf("Response status code: %v\n", resp.Status)

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	journal[url] = resp.StatusCode

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	log.Printf("success execute handler")

	return
}

package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type cachedResponse struct {
	Body        []byte
	ContentType string
	StatusCode  int
	ExpiresAt   time.Time
}

const (
	https       = "https://"
	contentType = "application/json"
	cacheTime   = 1 * time.Minute
)

var (
	journal = make(map[string]int)
	cache   = make(map[string]cachedResponse)
	mu      sync.Mutex

	URLs        []bannedURL
	banedURLMap = make(map[string]bool)
)

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
	_, ok := banedURLMap[url]
	if ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("banned site!"))
		log.Printf("Request for banned site %s", url)

		return
	}

	mu.Lock()
	cachedResp, exists := cache[url]
	mu.Unlock()

	if exists && time.Now().Before(cachedResp.ExpiresAt) {
		log.Printf("Cached response found for %s", url)
		w.Header().Set("Content-Type", cachedResp.ContentType)
		w.WriteHeader(http.StatusNotModified)
		w.Write(cachedResp.Body)

		return
	}

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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		mu.Lock()
		cache[url] = cachedResponse{
			Body:        body,
			ContentType: resp.Header.Get("Content-Type"),
			StatusCode:  resp.StatusCode,
			ExpiresAt:   time.Now().Add(cacheTime),
		}
		mu.Unlock()

		w.Write(body)

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	mu.Lock()
	cache[url] = cachedResponse{
		Body:        body,
		ContentType: resp.Header.Get("Content-Type"),
		StatusCode:  resp.StatusCode,
		ExpiresAt:   time.Now().Add(cacheTime),
	}
	mu.Unlock()

	w.Write(body)

	log.Printf("success execute handler")

	return
}

type bannedURL struct {
	U string `json:"u"`
}

func BuildBannedURL() {
	file, err := os.Open("banned_sites.json")
	if err != nil {
		log.Fatalf("Error l openning file:%v", err)
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&URLs)
	if err != nil {
		log.Fatalf("Error while decode json file:%v", err)
	}

	for _, url := range URLs {
		banedURLMap[url.U] = true
	}
}

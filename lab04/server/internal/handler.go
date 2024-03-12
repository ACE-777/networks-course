package internal

import (
	"io"
	"log"
	"net/http"
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

	//client := &http.Client{}
	//_, err := client.Get(destinationURL.String())
	//if err != nil {
	//	fmt.Println("errr:", err)
	//}

	resp, err := http.DefaultTransport.RoundTrip(proxyRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)

		return
	}

	//resp, err := client.Do(proxyRequest)
	//if err != nil {
	//	log.Printf("Error proxying request: %v\n", err)
	//	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//	return
	//}
	//
	defer resp.Body.Close()

	log.Printf("Response status code: %d\n", resp.StatusCode)

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("succed execute handler")
}

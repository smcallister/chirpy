package main

import (
	"net/http"
)

func healthzHandler(res http.ResponseWriter, req *http.Request) {
	// Write the response.
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}

func main() {
	// Set up the handler.
	handler := http.NewServeMux()
	handler.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	handler.Handle("/assets/", http.FileServer(http.Dir(".")))
	handler.HandleFunc("/healthz", healthzHandler)

	// Create the server.	
	server := http.Server{Handler: handler, Addr: ":8080"}

	// Listen for requests.
	server.ListenAndServe()
}
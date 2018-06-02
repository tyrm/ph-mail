package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func StartWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/envelope", handleEnvelope).Methods("GET")
	r.HandleFunc("/envelope", handleNotImplemented)

	// 404 handler
	r.PathPrefix("/").HandlerFunc(catchAllHandler)

	http.ListenAndServe("localhost:8080", r)
}

func handleEnvelope(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	fmt.Fprint(response, "{\"cat\":\"meow\"}")

	return
}

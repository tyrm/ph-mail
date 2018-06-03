package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../models"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

func StartWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/envelope/{messageId}", handleEnvelope).Methods("GET")
	r.HandleFunc("/envelope/{messageId}", handleNotImplemented)
	r.HandleFunc("/envelope", handleEnvelopes).Methods("GET")
	r.HandleFunc("/envelope", handleNotImplemented)

	// 404 handler
	r.PathPrefix("/").HandlerFunc(catchAllHandler)

	http.ListenAndServe(":8080", r)
}

func handleEnvelopes(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	fmt.Fprint(response, "{\"cat\":\"meow\"}")

	return
}

func handleEnvelope(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(request)
	envelope, error := models.GetEnvelopeByMsgID(vars["messageId"])

	if error == gorm.ErrRecordNotFound {
		makeErrorResponse(response, 404, vars["messageId"], 0)
		return
	} else if error != nil {
		makeErrorResponse(response, 400, error.Error(), 0)
		return
	}

	b, _ := json.Marshal(envelope)
	fmt.Fprintf(response, "%s", b)

	return
}
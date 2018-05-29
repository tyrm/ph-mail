package web

import (
	"../env"
	"net/http"
)

func AddHttpHandlers() {
	http.HandleFunc("/event/v1", handleFilterList)
}

func (myEnv *env.Env) handleFilterList(response http.ResponseWriter, request *http.Request) {

}
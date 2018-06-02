package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ErrorMessage struct {
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Status string `json:"status,omitempty"`
	Code   string `json:"code,omitempty"`
}

type ErrorResponse struct {
	Error  *ErrorMessage `json:"error"`
}

func catchAllHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	makeErrorResponse(response, 404, request.URL.Path, 0)
	return
}

func makeErrorResponse(response http.ResponseWriter, status int, detail string, code int) {
	var codeTitle map[int]string
	codeTitle = make(map[int]string)
	codeTitle[1] = "Malformed JSON Body"
	codeTitle[2201] = "Missing Required Attribute"
	codeTitle[2202] = "Requested Relationship Not Found"

	var statusTitle map[int]string
	statusTitle = make(map[int]string)
	statusTitle[400] = "Bad Request"
	statusTitle[401] = "Unauthorized"
	statusTitle[404] = "Not Found"
	statusTitle[405] = "Method Not Allowed"
	statusTitle[406] = "Not Acceptable"
	statusTitle[409] = "Conflict"
	statusTitle[415] = "Unsupported Media Type"
	statusTitle[422] = "Unprocessable Entity"
	statusTitle[500] = "Internal Server Error"

	var title string
	var statusStr string = strconv.Itoa(status)
	var codeStr string

	// Get Title
	if code == 0 { // code 0 means no code
		title = statusTitle[status]
	} else {
		title = codeTitle[code]
		codeStr = strconv.Itoa(code)
	}

	// Send Response
	response.WriteHeader(status)

	m := ErrorResponse{
		Error: &ErrorMessage{
			Title:  title,
			Detail: detail,
			Status: statusStr,
			Code:   codeStr,
		},
	}
	b, _ := json.Marshal(m)
	fmt.Fprintf(response, "%s", b)

	return
}

func handleNotImplemented(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	makeErrorResponse(response, 405, request.Method, 0)
	return
}
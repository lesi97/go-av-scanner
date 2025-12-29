package utils

import (
	"encoding/json"
	"net/http"

	"github.com/lesi97/go-av-scanner/internal/scanner"
)

type apiResponse struct {
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
}

/*
Success sends a successful JSON response with the given data

Response format:
	{ message: <data>, error: null }
Notes:
	- If `data` is a struct, all fields that should appear in the JSON must be public or they will be omitted

Example:
func test(w http.ResponseWriter, r *http.Request) {
	utils.Success(w, http.StatusOK, "message")
}
*/
func Success(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	response := apiResponse{
		Message: data,
		Error:   nil,
	}
	json.NewEncoder(w).Encode(response)
}

/*
Error sends a failed JSON response with the given error

Response format:
	{ message: null, error: <message> }
Notes:
	- You can only pass a type of Error or string
Example:
func test(w http.ResponseWriter, r *http.Request) {
	utils.Error(w, http.StatusInternalServerError, "An unexpected error has occurred!")
}
*/
func Error(w http.ResponseWriter, status int, errorMessage interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	var output interface{}

	switch v := errorMessage.(type) {
	case nil:
		output = "An unexpected error has occurred!"
	case error:
		output = v.Error()
	case string:
		output = v
	case scanner.Result:
		output = v
	case *scanner.Result:
		output = v
	case scanner.ScanError:
		output = v.Result
	default:
		output = "An unexpected error has occurred!"
	}

	response := apiResponse{
		Message: nil,
		Error:   output,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func TextResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(message))
}

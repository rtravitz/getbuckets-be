package handler

import (
	"encoding/json"
	"net/http"
)

//ErrorResponse exists to serialize an error message
type ErrorResponse struct {
	Error string `json:"error"`
}

func respond(w http.ResponseWriter, data interface{}, statusCode int) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func respondError(w http.ResponseWriter, err error) error {
	er := ErrorResponse{
		Error: err.Error(),
	}
	return respond(w, er, http.StatusInternalServerError)
}

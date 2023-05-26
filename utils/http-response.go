package utils

import (
	"log"
	"net/http"
)

const (
	headerContentType = "Content-Type"
	mimeJson          = "application/json"
	mimeText          = "plain/text"
)

// HttpOkResponse writes a 200 OK response with the given byte array as the body.
// The Content-Type header is set to application/json.
func HttpOkResponse(w http.ResponseWriter, b []byte) {
	w.Header().Add(headerContentType, mimeJson)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// HttpCreatedResponse writes a 201 Created response with the given byte array as the body.
// The Content-Type header is set to application/json. It logs the given message.
func HttpCreatedResponse(w http.ResponseWriter, b []byte, location string, msg string) {
	log.Print(msg)
	w.Header().Add(headerContentType, mimeJson)
	w.Header().Add("Location", location)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

// HttpNoContentResponse writes a 404 Not Found response with the given message as the body.
// The Content-Type header is set to plain/text. It logs the given message.
func HttpNotFoundResponse(w http.ResponseWriter, msg string) {
	log.Print(msg)
	w.Header().Add(headerContentType, mimeText)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(msg))
}

// HttpBadRequestResponse writes a 400 Bad Request response with the given message as the body.
// The Content-Type header is set to plain/text. It logs the given message.
func HttpBadRequestResponse(w http.ResponseWriter, msg string) {
	log.Print(msg)
	w.Header().Add(headerContentType, mimeText)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
}

// HttpInternalServerErrorResponse writes a 500 Internal Server Error response with the given message as the body.
// The Content-Type header is set to plain/text. It logs the given message.
func HttpInternalServerErrorResponse(w http.ResponseWriter, msg string) {
	log.Print(msg)
	w.Header().Add(headerContentType, mimeText)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}

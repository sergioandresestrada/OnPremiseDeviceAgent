package utils

import "net/http"

func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func OKRequest(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func ServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

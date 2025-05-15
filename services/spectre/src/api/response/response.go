package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

// Ok wrap data to json ok response
func Ok(w http.ResponseWriter, data interface{}) {
	r := Response{
		Success: true,
		Data:    data,
	}
	json.NewEncoder(w).Encode(r)
}

func ErrFailedToRetrieveLetters(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	r := Response{
		Success: false,
		Error:   "failed to retrieve letters",
	}
	json.NewEncoder(w).Encode(r)
}

func ErrInvalidID(w http.ResponseWriter, sid string) {
	w.WriteHeader(http.StatusBadRequest)
	r := Response{
		Success: false,
		Error:   "invalid id: " + sid,
	}
	json.NewEncoder(w).Encode(r)
}

func ErrNotFound(w http.ResponseWriter, sid string) {
	w.WriteHeader(http.StatusNotFound)
	r := Response{
		Success: false,
		Error:   "letter with id " + sid + " not found",
	}
	json.NewEncoder(w).Encode(r)
}

func ErrCannotGetWithID(w http.ResponseWriter, sid string) {
	w.WriteHeader(http.StatusInternalServerError)
	r := Response{
		Success: false,
		Error:   "cannot get with id " + sid,
	}
	json.NewEncoder(w).Encode(r)
}

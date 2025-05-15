package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Content interface{} `json:"content"`
	Error   interface{} `json:"error"`
}

// Ok wrap data to json ok response
func Ok(w http.ResponseWriter, content interface{}) {
	r := Response{
		Success: true,
		Content: content,
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

func ErrCannotDeleteWithID(w http.ResponseWriter, sid string) {
	w.WriteHeader(http.StatusInternalServerError)
	r := Response{
		Success: false,
		Error:   "cannot delete with id " + sid,
	}
	json.NewEncoder(w).Encode(r)
}

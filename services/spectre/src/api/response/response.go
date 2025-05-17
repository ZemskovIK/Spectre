package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Content interface{} `json:"content"`
	Error   interface{} `json:"error"`
}

func NewResponse(content interface{}, err interface{}) Response {
	return Response{
		Content: content,
		Error:   err,
	}
}

// Ok wrap data to json ok response
func Ok(w http.ResponseWriter, content interface{}) {
	r := Response{
		Content: content,
	}
	json.NewEncoder(w).Encode(r)
}

func ErrFailedToRetrieveLetters(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	r := NewResponse(nil, "failed to retrieve letters")
	json.NewEncoder(w).Encode(r)
}

func ErrInvalidID(w http.ResponseWriter, sid string) {
	w.WriteHeader(http.StatusBadRequest)
	r := NewResponse(nil, "invalid id: "+sid)
	json.NewEncoder(w).Encode(r)
}

func ErrNotFound(w http.ResponseWriter, sid string) {
	w.WriteHeader(http.StatusNotFound)
	r := NewResponse(nil, "letter with id "+sid+" not found")
	json.NewEncoder(w).Encode(r)
}

func ErrCannotGetWithID(w http.ResponseWriter, sid string) {
	w.WriteHeader(http.StatusInternalServerError)
	r := NewResponse(nil, "cannot get with id "+sid)
	json.NewEncoder(w).Encode(r)
}

func ErrCannotDeleteWithID(w http.ResponseWriter, sid string) {
	w.WriteHeader(http.StatusInternalServerError)
	r := NewResponse(nil, "cannot delete with id "+sid)
	json.NewEncoder(w).Encode(r)
}

func ErrInvalidRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	r := NewResponse(nil, msg)
	json.NewEncoder(w).Encode(r)
}

func ErrCannotSave(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	r := NewResponse(nil, "cannot save!")
	json.NewEncoder(w).Encode(r)
}

func ErrCannotUpdate(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	r := NewResponse(nil, "cannot update!")
	json.NewEncoder(w).Encode(r)
}

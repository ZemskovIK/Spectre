package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Content interface{} `json:"content"`
	Error   interface{} `json:"error"`

	IV    string `json:"iv"`
	HMAC  string `json:"hmac"`
	Nonce string `json:"nonce"`
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

func JWTOk(w http.ResponseWriter, token string) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
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

func ErrCannotGetCredsFromJSON(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	r := NewResponse(nil, "cannot get creds from json!")
	json.NewEncoder(w).Encode(r)
}

func ErrCannotSignIn(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	r := NewResponse(nil, "cannot sign in with jwt!")
	json.NewEncoder(w).Encode(r)
}

func ErrCannotGetUserByLogin(w http.ResponseWriter, login string) {
	w.WriteHeader(http.StatusInternalServerError)
	r := NewResponse(nil, "cannot get user by login: "+login)
	json.NewEncoder(w).Encode(r)
}

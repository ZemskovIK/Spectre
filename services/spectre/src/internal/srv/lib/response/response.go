package response

// ! TODO : think about copy-paste

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response is a standard API response structure.
type ResponseWithContent struct {
	Content interface{} `json:"content"`
	Error   interface{} `json:"error"`

	IV    string `json:"iv"`
	HMAC  string `json:"hmac"`
	Nonce string `json:"nonce"`
}

type ECDHResponse struct {
	Key string `json:"key"`
}

var (
	EmptyWithContent = ResponseWithContent{}
	EmptyEDCH        = ECDHResponse{}
)

// NewWithContent creates a new ResponseWithContent object with err as second arg (nil if not).
func NewWithContent(content interface{}, err interface{}) ResponseWithContent {
	return ResponseWithContent{
		Content: content,
		Error:   err,
	}
}

// OkWithContent wraps data in a successful JSON response.
func OkWithContent(w http.ResponseWriter, content interface{}) {
	r := ResponseWithContent{
		Content: content,
	}
	json.NewEncoder(w).Encode(r)
}

// JWTOk sends a JSON response with a JWT token.
func JWTOk(w http.ResponseWriter, token string) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// OkWithResponse encode json r ResponseWithContent in w ResponseWrites.
func OkWithResponse(w http.ResponseWriter, r ResponseWithContent) {
	json.NewEncoder(w).Encode(r)
}

// OkWithECDHKey encode json r ECDHResponse in w ResponseWriter.
func OkWithECDHKey(w http.ResponseWriter, r ECDHResponse) {
	json.NewEncoder(w).Encode(r)
}

// OkEmpty writes in w ResponseWriter only 204.
func OkEmpty(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// ErrWithContent sends an error response with a given status and message.
func ErrWithContent(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	r := NewWithContent(nil, msg)
	json.NewEncoder(w).Encode(r)
}

func ErrWithECDH(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadGateway)
	r := map[string]string{
		"error": msg,
	}
	json.NewEncoder(w).Encode(r)
}

// ErrInvalidID sends a 400 error for invalid ID.
func ErrInvalidID(w http.ResponseWriter, sid string) {
	ErrWithContent(w, http.StatusBadRequest, "invalid id: "+sid)
}

// ErrNotFound sends a 404 error when a letter is not found.
func ErrNotFound(w http.ResponseWriter, sid string) {
	ErrWithContent(w, http.StatusNotFound, "letter with id "+sid+" not found")
}

// ErrCannotGetWithID sends a 500 error when unable to get by ID.
func ErrCannotGetWithID(w http.ResponseWriter, sid string) {
	ErrWithContent(w, http.StatusInternalServerError, "cannot get with id "+sid)
}

// ErrCannotDeleteWithID sends a 500 error when unable to delete by ID.
func ErrCannotDeleteWithID(w http.ResponseWriter, sid string) {
	ErrWithContent(w, http.StatusInternalServerError, "cannot delete with id "+sid)
}

// ErrInvalidRequest sends a 400 error for invalid requests.
func ErrInvalidRequest(w http.ResponseWriter, msg string) {
	ErrWithContent(w, http.StatusBadRequest, msg)
}

// ErrCannotSave sends a 422 error when unable to save.
func ErrCannotSave(w http.ResponseWriter) {
	ErrWithContent(w, http.StatusUnprocessableEntity, "cannot save!")
}

// ErrCannotUpdate sends a 500 error when unable to update.
func ErrCannotUpdate(w http.ResponseWriter) {
	ErrWithContent(w, http.StatusInternalServerError, "cannot update!")
}

// ErrCannotGetCredsFromJSON sends a 500 error when unable to get credentials from JSON.
func ErrCannotGetCredsFromJSON(w http.ResponseWriter) {
	ErrWithContent(w, http.StatusInternalServerError, "cannot get creds from json!")
}

// ErrCannotSignIn sends a 500 error when sign-in fails.
func ErrCannotSignIn(w http.ResponseWriter) {
	ErrWithContent(w, http.StatusInternalServerError, "cannot sign in with jwt!")
}

// ErrCannotGetUserByLogin sends a 500 error when unable to get user by login.
func ErrCannotGetUserByLogin(w http.ResponseWriter, login string) {
	ErrWithContent(w, http.StatusInternalServerError, "cannot get user by login: "+login)
}

// ErrBlockedToGet sends a 403 error when access is blocked due to insufficient level.
func ErrBlockedToGet(w http.ResponseWriter, usrAL, neededAL int) {
	ErrWithContent(
		w,
		http.StatusForbidden,
		fmt.Sprintf("blocked to get letter with access level %d with your level %d", neededAL, usrAL),
	)
}

// ErrYouArntAdmin sends a 403 error when user is not admin.
func ErrYouArntAdmin(w http.ResponseWriter) {
	ErrWithContent(w, http.StatusForbidden, "you have not access to admin panel!")
}

// ErrCannotRetrieveLetters sends a 500 error when letters cannot be retrieved.
func ErrCannotRetrieveLetters(w http.ResponseWriter) {
	ErrWithContent(w, http.StatusInternalServerError, "cannot retrieve letters")
}

// ErrCannotRetrieveUsers sends a 500 error when users cannot be retrieved.
func ErrCannotRetrieveUsers(w http.ResponseWriter) {
	ErrWithContent(w, http.StatusInternalServerError, "cannot retrieve users")
}

// ErrCannotGetB64Strings sends a 500 error when data cannot be converted to base64.
func ErrCannotGetB64Strings(w http.ResponseWriter) {
	ErrWithContent(w, http.StatusInternalServerError, "cannot convert data to base64")
}

func ErrCannotFetchFromB64(w http.ResponseWriter) {
	ErrWithContent(w, http.StatusInternalServerError, "cannot fetch from base64")
}

// ErrCannotEncryptData sends a 502 error when encryption service is unavailable or failed.
func ErrCannotEncryptData(w http.ResponseWriter) {
	ErrWithContent(w, http.StatusBadGateway, "cannot encrypt data!")
}

// ErrCannotDecryptData sends a 502 error when decryption service is unavailable or failed.
func ErrCannotDecryptData(w http.ResponseWriter) {
	ErrWithContent(w, http.StatusBadGateway, "cannot decrypt data!")
}

// ErrCannotECDHGetK sends a 502 error when ECDH fails
func ErrCannotECDHGetK(w http.ResponseWriter) {
	ErrWithECDH(w, "cannot ECDH get K")
}

// ErrCannotECDHSetA sends a 502 error when ECDH fails
func ErrCannotECDHSetA(w http.ResponseWriter) {
	ErrWithECDH(w, "cannot ECDH set A")
}

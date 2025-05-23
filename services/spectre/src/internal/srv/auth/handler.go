package auth

import (
	"encoding/json"
	"net/http"
	"spectre/internal/models"
	"spectre/internal/srv/response"
	"spectre/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const GLOC = "src/internal/server/auth/handler.go/" // for logging

type authStore interface {
	GetUserByLogin(login string) (models.User, error)
}

type authHandler struct {
	st  authStore
	log *logger.Logger
}

func NewAuthHandler(
	s authStore, log *logger.Logger,
) *authHandler {
	return &authHandler{
		st:  s,
		log: log,
	}
}

// Login authenticates a user and returns a JWT token.
func (h *authHandler) Login(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC + "Login()"
	h.log.Infof("%s: handler called", loc)
	h.log.Debugf("%s: request: %+v", loc, r)

	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		h.log.Errorf("%s: failed to decode credentials from request body: %v", loc, err)
		response.ErrCannotGetCredsFromJSON(w)
		return
	}
	// ! TODO : decrypt
	h.log.Debugf("%s: credentials received for login: %s", loc, creds.Login)

	user, err := h.st.GetUserByLogin(creds.Login)
	if err != nil {
		h.log.Errorf("%s: failed to get user by login '%s': %v", loc, creds.Login, err)
		response.ErrCannotGetUserByLogin(w, creds.Login)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		user.PHash,
		[]byte(creds.Password),
	); err != nil {
		h.log.Warnf("%s: invalid password for user: %s", loc, creds.Login)
		// response.ErrCannotSignIn(w) // ! WARN
		// return
	}

	h.log.Infof("%s: user %s authenticated, generating JWT", loc, creds.Login)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.AccessLevel,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	stoken, err := token.SignedString([]byte("test_secret")) // ! TODO
	if err != nil {
		h.log.Errorf("%s: failed to sign JWT for user '%s': %v", loc, creds.Login, err)
		response.ErrCannotSignIn(w)
		return
	}

	h.log.Debugf("%s: JWT generated for user: %s, token: %s", loc, creds.Login, stoken)
	response.JWTOk(w, stoken)
}

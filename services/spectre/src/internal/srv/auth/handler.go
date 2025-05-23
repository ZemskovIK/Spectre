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

const GLOC = "src/internal/server/auth/handler.go/"

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

func (h *authHandler) Login(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC + "Login()"
	_ = loc

	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		response.ErrCannotGetCredsFromJSON(w)
		return
	}
	// ! TODO : decrypt
	user, err := h.st.GetUserByLogin(creds.Login)
	if err != nil {
		h.log.Errorf("%s: failed to get user by login: %v", loc, err)
		response.ErrCannotGetUserByLogin(w, creds.Login)
		return
	}
	if err := bcrypt.CompareHashAndPassword( // ! TODO
		user.PHash,
		[]byte(creds.Password),
	); err != nil {

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // ! WARN ES256
		"sub":  user.ID,
		"role": user.AccessLevel,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	stoken, err := token.SignedString([]byte("test_secret")) // ! TODO
	if err != nil {
		response.ErrCannotSignIn(w)
		return
	}

	response.JWTOk(w, stoken)
}

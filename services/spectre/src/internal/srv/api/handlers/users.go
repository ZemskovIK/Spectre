package handlers

import (
	"net/http"
	"spectre/internal/lib"
	"spectre/internal/models"
	"spectre/internal/srv/response"
	"spectre/pkg/logger"
)

const GLOC_USRS = "src/internal/api/handlers/users.go"

type usersStore interface {
	GetUserByLogin(login string) (models.User, error)
	SaveUser(usr models.User) error
	DeleteUser(id int) error
	UpdateUser(usr models.User) error
	GetAllUsers() ([]models.User, error)
}

type usersHandler struct {
	st  usersStore
	log *logger.Logger
}

func NewUsersHandler(
	s usersStore, log *logger.Logger,
) *usersHandler {
	return &usersHandler{
		st:  s,
		log: log,
	}
}

func (h *usersHandler) GetAll(
	w http.ResponseWriter, r *http.Request,
) {
	usrAccess, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		response.ErrFailedToRetrieveUsers(w)
		return
	}

	if usrAccess < lib.ADMIN_ALEVEL {
		response.ErrBlockedToGet(w, usrAccess, lib.ADMIN_ALEVEL)
		return
	}

	users, err := h.st.GetAllUsers()
	if err != nil {
		response.ErrFailedToRetrieveUsers(w)
		return
	}

	if len(users) == 0 {
		response.Ok(w, []interface{}{})
		return
	}

	// ! TODO : encrypt

	response.Ok(w, users)
}

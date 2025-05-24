package handlers

import (
	"net/http"
	"spectre/internal/models"
	"spectre/internal/srv/lib"
	"spectre/internal/srv/lib/response"
	"spectre/pkg/logger"
)

const GLOC_USRS = "src/internal/api/handlers/users.go" // for logging

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

// GetAll returns all users (admin only).
func (h *usersHandler) GetAll(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_USRS + "GetAll()"
	h.log.Infof("%s: handler called", loc)
	h.log.Debugf("%s: request: %+v", loc, r)

	usrAccess, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		h.log.Errorf("%s: failed to fetch access level from context", loc)
		response.ErrCannotRetrieveUsers(w)
		return
	}

	if usrAccess < lib.ADMIN_ALEVEL {
		h.log.Warnf("%s: blocked to get users, access: %d, required: %d", loc, usrAccess, lib.ADMIN_ALEVEL)
		response.ErrBlockedToGet(w, usrAccess, lib.ADMIN_ALEVEL)
		return
	}

	users, err := h.st.GetAllUsers()
	if err != nil {
		h.log.Errorf("%s: failed to retrieve users from storage: %v", loc, err)
		response.ErrCannotRetrieveUsers(w)
		return
	}

	if len(users) == 0 {
		h.log.Warnf("%s: no users found", loc)
		response.Ok(w, []interface{}{})
		return
	}

	h.log.Debugf("%s: successfully retrieved %d users", loc, len(users))
	response.Ok(w, users)
}

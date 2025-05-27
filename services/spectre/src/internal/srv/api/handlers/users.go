package handlers

import (
	"encoding/json"
	"net/http"
	"spectre/internal/models"
	"spectre/internal/srv/api"
	"spectre/internal/srv/lib"
	"spectre/internal/srv/lib/response"
	st "spectre/internal/storage"
	"spectre/pkg/logger"
	"strconv"
)

const GLOC_USRS = "src/internal/api/handlers/users.go/" // for logging

type usersStore interface {
	GetUserByLogin(login string) (models.User, error)
	GetUserByID(id int) (models.User, error)
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

	// ! TODO : copy-paste
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

func (h *usersHandler) GetOne(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_USRS + "GetOne()"

	// ! TODO : copy-paste
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

	sid, id, err := lib.GetID(api.USER_POINT, r.RequestURI)
	if err != nil {
		h.log.Errorf("%s: invalid id '%s' in URI '%s': %v", loc, sid, r.RequestURI, err)
		response.ErrInvalidID(w, sid)
		return
	}

	user, err := h.st.GetUserByID(id)
	if err != nil {
		response.ErrCannotGetWithID(w, sid)
		return
	}

	// ! TODO : encrypt

	response.Ok(w, user)
}

func (h *usersHandler) Delete(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_USRS + "Delete()"

	// ! TODO : copy-paste
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

	sid, id, err := lib.GetID(api.USER_POINT, r.RequestURI)
	if err != nil {
		h.log.Errorf("%s: invalid id '%s' in URI '%s': %v", loc, sid, r.RequestURI, err)
		response.ErrInvalidID(w, sid)
		return
	}

	if err := h.st.DeleteUser(id); err != nil {
		response.ErrCannotDeleteWithID(w, sid)
		return
	}

	// ! TODO : encrypt

	response.Ok(w, nil)
}

func (h *usersHandler) Add(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_USRS + "Add()"

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

	type user struct {
		models.User
		Password string `json:"password"`
	}
	var usr user
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		h.log.Errorf("%s: failed to decode JSON body: %v", loc, err)
		response.ErrInvalidRequest(w, "invalid JSON")
		return
	}
	usr.PHash = []byte(usr.Password) // ! TODO
	h.log.Debugf("%s: decoded usr: %+v", loc, usr)

	if usr.Login == "" {
		response.ErrInvalidRequest(w, "login cannot be empty!")
		return
	}
	if usr.Password == "" {
		response.ErrInvalidRequest(w, "password cannot be empty!")
		return
	}
	if usr.AccessLevel < 1 || usrAccess > lib.ADMIN_ALEVEL {
		response.ErrInvalidRequest(w, "invalid access_level "+strconv.Itoa(usr.AccessLevel))
		return
	}

	// ! TODO : hash pass

	if err := h.st.SaveUser(usr.User); err != nil {
		response.ErrCannotSave(w)
		return
	}

	response.Ok(w, nil)
}

func (h *usersHandler) Update(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_USRS + "Update()"

	usrAccessLevel, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		h.log.Errorf("%s: failed to fetch access level from context", loc)
		response.ErrCannotUpdate(w)
		return
	}
	if usrAccessLevel < lib.ADMIN_ALEVEL {
		h.log.Warnf("%s: blocked to get users, access: %d, required: %d", loc, usrAccessLevel, lib.ADMIN_ALEVEL)
		response.ErrYouArntAdmin(w)
		return
	}

	sid, id, err := lib.GetID(api.USER_POINT, r.RequestURI)
	if err != nil {
		h.log.Errorf("%s: invalid id '%s' in URI '%s': %v", loc, sid, r.RequestURI, err)
		response.ErrInvalidID(w, sid)
		return
	}

	type user struct {
		models.User
		Password string `json:"password"`
	}
	var usr user
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		h.log.Errorf("%s: failed to decode JSON body: %v", loc, err)
		response.ErrInvalidRequest(w, "invalid JSON")
		return
	}
	usr.ID = id
	usr.PHash = []byte(usr.Password)

	// ! TODO validation func
	if usr.Login == "" {
		response.ErrInvalidRequest(w, "login cannot be empty!")
		return
	}
	if usr.Password == "" {
		response.ErrInvalidRequest(w, "password cannot be empty!")
		return
	}
	if usr.AccessLevel < 1 || usrAccessLevel > lib.ADMIN_ALEVEL {
		response.ErrInvalidRequest(w, "invalid access_level "+strconv.Itoa(usr.AccessLevel))
		return
	}

	if err := h.st.UpdateUser(usr.User); err != nil {
		if err.Error() == st.ErrLetterNotFound(id).Error() {
			h.log.Warnf("%s: user with id not found: %v", loc, err)
			response.ErrNotFound(w, sid)
		} else {
			h.log.Errorf("%s: error when updating user: %v", loc, err)
			response.ErrCannotUpdate(w)
		}
		return
	}

	response.Ok(w, nil)
}

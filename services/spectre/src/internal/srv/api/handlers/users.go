package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"spectre/internal/models"
	"spectre/internal/srv/api"
	"spectre/internal/srv/lib"
	"spectre/internal/srv/lib/response"
	"spectre/internal/srv/proxy"
	"spectre/pkg/logger"

	"golang.org/x/crypto/bcrypt"
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
	crypto *proxy.CryptoClient
	st     usersStore
	log    *logger.Logger
}

func NewUsersHandler(
	s usersStore, log *logger.Logger, cr *proxy.CryptoClient,
) *usersHandler {
	return &usersHandler{
		crypto: cr,
		st:     s,
		log:    log,
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
		response.OkWithContent(w, []interface{}{})
		return
	}

	h.log.Debugf("%s: successfully retrieved %d users", loc, len(users))
	response.OkWithContent(w, users)
}

func (h *usersHandler) ECDHGetK(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_USRS + "ECDHGetK()"
	h.log.Debugf("%s: handler called", loc)

	resp, err := h.crypto.GetK(r)
	if err != nil {
		h.log.Errorf("%s: error when get k from proxy: %v", loc, err)
		response.ErrCannotECDHGetK(w)
		return
	}

	response.OkWithECDHKey(w, resp)
}

func (h *usersHandler) ECDHSetA(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_USRS + "ECDHSetA()"
	h.log.Debugf("%s: handler called", loc)

	if err := h.crypto.SetA(r); err != nil {
		h.log.Errorf("%s: error when set a from proxy: %v", loc, err)
		response.ErrCannotECDHSetA(w)
		return
	}

	response.OkEmpty(w)
}

func (h *usersHandler) Create(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_USRS + "Create()"
	h.log.Infof("%s: handler called", loc)

	usrAccess, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		h.log.Errorf("%s: failed to fetch access level from context", loc)
		response.ErrWithContent(w, http.StatusInternalServerError, "cannot get access level")
		return
	}
	if usrAccess < lib.ADMIN_ALEVEL {
		h.log.Warnf("%s: blocked, access: %d, required: %d", loc, usrAccess, lib.ADMIN_ALEVEL)
		response.ErrWithContent(w, http.StatusForbidden, "admin access required")
		return
	}

	var encryptedReq map[string]interface{}
	if err := lib.ReadJSON(r, &encryptedReq); err != nil {
		h.log.Errorf("%s: failed to read request body: %v", loc, err)
		response.ErrWithContent(w, http.StatusBadRequest, "invalid request body")
		return
	}

	h.log.Debugf("%s: received encrypted request: %+v", loc, encryptedReq)

	decrypted, err := h.crypto.DecryptData(encryptedReq, r.Host)
	if err != nil {
		h.log.Errorf("%s: failed to decrypt data: %v", loc, err)
		response.ErrCannotDecryptData(w)
		return
	}

	contentList, ok := decrypted.Content.([]interface{})
	if !ok || len(contentList) == 0 {
		h.log.Errorf("%s: unexpected decrypted content format: %+v", loc, decrypted.Content)
		response.ErrWithContent(w, http.StatusBadRequest, "invalid decrypted content")
		return
	}

	userJSONStr, ok := contentList[0].(string)
	if !ok {
		h.log.Errorf("%s: decrypted content[0] is not a string", loc)
		response.ErrWithContent(w, http.StatusBadRequest, "invalid user data format")
		return
	}

	userJSONBytes, err := base64.StdEncoding.DecodeString(userJSONStr)
	if err != nil {
		h.log.Errorf("%s: failed to base64 decode user data: %v", loc, err)
		response.ErrWithContent(w, http.StatusBadRequest, "invalid base64 user data")
		return
	}

	var req map[string]interface{}
	if err := lib.ReadJSONFromBytes(userJSONBytes, &req); err != nil {
		h.log.Errorf("%s: failed to parse user JSON: %v", loc, err)
		response.ErrWithContent(w, http.StatusBadRequest, "invalid user json")
		return
	}

	h.log.Debugf("%s: parsed user data: %+v", loc, req)

	login, _ := req["login"].(string)
	password, _ := req["password"].(string)
	accessLevel, _ := req["access_level"].(float64)

	if login == "" || password == "" {
		response.ErrWithContent(w, http.StatusBadRequest, "login and password required")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		h.log.Errorf("%s: failed to hash password: %v", loc, err)
		response.ErrWithContent(w, http.StatusInternalServerError, "cannot hash password")
		return
	}

	newUser := models.User{
		Login:       login,
		PHash:       hashedPassword,
		AccessLevel: int(accessLevel),
	}

	if err := h.st.SaveUser(newUser); err != nil {
		h.log.Errorf("%s: failed to save user: %v", loc, err)
		response.ErrWithContent(w, http.StatusInternalServerError, "cannot save user")
		return
	}

	h.log.Infof("%s: user created successfully: %s", loc, login)
	response.OkWithContent(w, "User created successfully")
}

// GetOne returns a single user by ID or login (admin only)
func (h *usersHandler) GetOne(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_USRS + "GetOne()"
	h.log.Infof("%s: handler called", loc)

	// Check access
	usrAccess, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		h.log.Errorf("%s: failed to fetch access level from context", loc)
		response.ErrCannotRetrieveUsers(w)
		return
	}
	if usrAccess < lib.ADMIN_ALEVEL {
		h.log.Warnf("%s: blocked to get user, access: %d, required: %d", loc, usrAccess, lib.ADMIN_ALEVEL)
		response.ErrBlockedToGet(w, usrAccess, lib.ADMIN_ALEVEL)
		return
	}

	// Get ID from URL
	idStr := r.URL.Path[len(api.USER_POINT):]
	if idStr == "" {
		h.log.Errorf("%s: empty user ID", loc)
		response.ErrInvalidID(w, idStr)
		return
	}

	var user models.User

	// Transform ID to int
	id, err := lib.ParseID(idStr)
	// Get user from DB
	if err == nil {
		user, err = h.st.GetUserByID(id)
		if err != nil {
			h.log.Errorf("%s: failed to get user by ID %d: %v", loc, id, err)
			response.ErrCannotGetWithID(w, idStr)
			return
		}
	} else {
		h.log.Debugf("%s: '%s' is not numeric, searching by login", loc, idStr)
		user, err = h.st.GetUserByLogin(idStr)
		if err != nil {
			h.log.Errorf("%s: failed to get user by login '%s': %v", loc, idStr, err)
			response.ErrWithContent(w, http.StatusNotFound, fmt.Sprintf("user with login '%s' not found", idStr))
			return
		}
	}

	h.log.Debugf("%s: successfully retrieved user: %+v", loc, user)
	response.OkWithContent(w, user)
}

func (h *usersHandler) Delete(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_USRS + "Delete()"
	h.log.Infof("%s: handler called", loc)

	usrAccess, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		h.log.Errorf("%s: failed to fetch access level from context", loc)
		response.ErrWithContent(w, http.StatusInternalServerError, "cannot get access level")
		return
	}
	if usrAccess < lib.ADMIN_ALEVEL {
		h.log.Warnf("%s: blocked, access: %d, required: %d", loc, usrAccess, lib.ADMIN_ALEVEL)
		response.ErrWithContent(w, http.StatusForbidden, "admin access required")
		return
	}

	idStr := r.URL.Path[len(api.USER_POINT):]
	if idStr == "" {
		response.ErrInvalidID(w, idStr)
		return
	}

	id, err := lib.ParseID(idStr)
	if err != nil {
		h.log.Errorf("%s: invalid user ID: %v", loc, err)
		response.ErrInvalidID(w, idStr)
		return
	}

	if err := h.st.DeleteUser(id); err != nil {
		h.log.Errorf("%s: failed to delete user %d: %v", loc, id, err)
		response.ErrWithContent(w, http.StatusInternalServerError, "cannot delete user")
		return
	}

	h.log.Infof("%s: user deleted successfully: %d", loc, id)
	response.OkWithContent(w, "User deleted successfully")
}

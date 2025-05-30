package handlers

import (
	"net/http"
	"spectre/internal/models"
	"spectre/internal/srv/api"
	"spectre/internal/srv/lib"
	"spectre/internal/srv/lib/response"
	"spectre/internal/srv/proxy"
	st "spectre/internal/storage"
	"spectre/pkg/logger"
)

const GLOC_LTS = "src/internal/api/handlers/letters.go/" // for logging

type lettersStore interface {
	GetAllLettersWithAccess(accessLevel int) ([]models.Letter, error)
	GetLetterByID(id int) (models.Letter, error)
	SaveLetter(letter models.Letter) error
	UpdateLetter(letter models.Letter) error
	DeleteLetter(id int) error
}

type lettersHandler struct {
	crypto *proxy.CryptoClient
	st     lettersStore
	log    *logger.Logger
}

// NewLettersHandler creates a new letters handler.
func NewLettersHandler(
	s lettersStore, log *logger.Logger, cr *proxy.CryptoClient,
) *lettersHandler {
	return &lettersHandler{
		st:     s,
		log:    log,
		crypto: cr,
	}
}

// GetAll returns all letters according to user's access level.
func (h *lettersHandler) GetAll(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_LTS + "GetAll()"
	h.log.Infof("%s: handler called", loc)
	h.log.Debugf("%s: request: %+v", loc, r)

	usrAccess, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		h.log.Errorf("%s: failed to fetch access level from context", loc)
		response.ErrCannotRetrieveLetters(w)
		return
	}

	letters, err := h.st.GetAllLettersWithAccess(usrAccess)
	if err != nil {
		h.log.Errorf("%s: failed to retrieve letters from storage: %v", loc, err)
		response.ErrCannotRetrieveLetters(w)
		return
	}

	if len(letters) == 0 {
		h.log.Warnf("%s: no letters found for access level %d", loc, usrAccess)
		response.OkWithContent(w, []interface{}{})
		return
	}

	h.log.Debugf("%s: found %d letters for access level %d", loc, len(letters), usrAccess)

	b64, err := lib.ToBase64Slice(letters)
	if err != nil {
		h.log.Errorf("%s: cannot convert to b64: %v", loc, err)
		response.ErrCannotGetB64Strings(w)
		return
	}
	resp, err := h.crypto.Encrypt(b64, r.Host)
	if err != nil {
		h.log.Errorf("%s: cannot enctypt: %v", loc, err)
		response.ErrCannotEncryptData(w)
		return
	}

	response.OkWithResponse(w, resp)
}

// GetOne returns a single letter by id according to user's access level.
func (h *lettersHandler) GetOne(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_LTS + "GetOne()"
	h.log.Infof("%s: handler called", loc)
	h.log.Debugf("%s: request: %+v", loc, r)

	usrAccess, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		h.log.Errorf("%s: failed to fetch access level from context", loc)
		response.ErrCannotRetrieveLetters(w)
		return
	}

	sid, id, err := lib.GetID(api.LETTER_POINT, r.RequestURI)
	if err != nil {
		h.log.Errorf("%s: invalid id '%s' in URI '%s': %v", loc, sid, r.RequestURI, err)
		response.ErrInvalidID(w, sid)
		return
	}

	letter, err := h.st.GetLetterByID(id)
	if err != nil {
		if err.Error() == st.ErrLetterNotFound(id).Error() {
			h.log.Warnf("%s: letter not found with id: %d", loc, id)
			response.ErrNotFound(w, sid)
		} else {
			h.log.Errorf("%s: failed to retrieve letter with id: %d, error: %v", loc, id, err)
			response.ErrCannotGetWithID(w, sid)
		}
		return
	}

	if usrAccess < letter.AccessLevel {
		h.log.Warnf("%s: user access level %d is lower than letter access level %d", loc, usrAccess, letter.AccessLevel)
		response.ErrBlockedToGet(w, usrAccess, letter.AccessLevel)
		return
	}
	h.log.Debugf("%s: successfully retrieved letter: %+v", loc, letter)

	b64, err := lib.ToBase64(letter)
	if err != nil {
		h.log.Errorf("%s: cannot converto to b64: %v", loc, err)
		response.ErrCannotGetB64Strings(w)
		return
	}
	resp, err := h.crypto.Encrypt([]string{b64}, r.Host)
	if err != nil {
		h.log.Errorf("%s: cannot enctypt: %v", loc, err)
		response.ErrCannotEncryptData(w)
		return
	}

	response.OkWithResponse(w, resp)
}

// Delete removes a letter by id.
func (h *lettersHandler) Delete(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_LTS + "Delete()"
	h.log.Infof("%s: handler called", loc)
	h.log.Debugf("%s: request: %+v", loc, r)

	usrAccess, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		h.log.Errorf("%s: failed to fetch access level from context", loc)
		response.ErrCannotRetrieveLetters(w)
		return
	}
	if usrAccess < lib.ADMIN_ALEVEL {
		h.log.Warnf("%s: user access level %d is not admin", loc, usrAccess)
		response.ErrYouArntAdmin(w)
		return
	}

	sid, id, err := lib.GetID(api.LETTER_POINT, r.RequestURI)
	if err != nil {
		h.log.Errorf("%s: invalid id '%s' in URI '%s': %v", loc, sid, r.RequestURI, err)
		response.ErrInvalidID(w, sid)
		return
	}

	h.log.Debugf("%s: deleting letter with id: %d", loc, id)
	if err := h.st.DeleteLetter(id); err != nil {
		if err.Error() == st.ErrLetterNotFound(id).Error() {
			h.log.Warnf("%s: letter not found with id: %d", loc, id)
			response.ErrNotFound(w, sid)
		} else {
			h.log.Errorf("%s: failed to delete letter with id: %d, error: %v", loc, id, err)
			response.ErrCannotDeleteWithID(w, sid)
		}
		return
	}

	h.log.Debugf("%s: letter with id %d deleted successfully", loc, id)
	response.OkWithContent(w, nil)
}

// Add creates a new letter (admin only).
func (h *lettersHandler) Add(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_LTS + "Add()"
	h.log.Infof("%s: handler called", loc)
	h.log.Debugf("%s: request: %+v", loc, r)

	usrAccessLevel, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		h.log.Errorf("%s: failed to fetch access level from context", loc)
		response.ErrCannotUpdate(w)
		return
	}
	if usrAccessLevel < lib.ADMIN_ALEVEL {
		h.log.Warnf("%s: user access level %d is not admin", loc, usrAccessLevel)
		response.ErrYouArntAdmin(w)
		return
	}

	resp, err := h.crypto.Decrypt(r)
	if err != nil {
		h.log.Errorf("%s: cannot decrypt %v", loc, err)
		response.ErrCannotDecryptData(w)
		return
	}
	var letter models.Letter
	if err := lib.FetchLetterFromB64(&letter, resp.Content); err != nil {
		h.log.Errorf("%s: cannot fetch from b64: %v", loc, err)
		response.ErrCannotFetchFromB64(w)
		return
	}

	msg, ok := lib.ValidateLetter(letter, h.log)
	if !ok {
		response.ErrInvalidRequest(w, msg)
		return
	}

	h.log.Debugf("%s: saving letter: %+v", loc, letter)
	if err := h.st.SaveLetter(letter); err != nil {
		h.log.Errorf("%s: failed to save letter: %v", loc, err)
		response.ErrCannotSave(w)
		return
	}

	h.log.Debugf("%s: letter saved successfully", loc)
	response.OkWithContent(w, nil)
}

// Update updates a letter by id (admin only).
func (h *lettersHandler) Update(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_LTS + "Update()"
	h.log.Infof("%s: handler called", loc)
	h.log.Debugf("%s: request: %+v", loc, r)

	usrAccessLevel, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		h.log.Errorf("%s: failed to fetch access level from context", loc)
		response.ErrCannotUpdate(w)
		return
	}
	if usrAccessLevel < lib.ADMIN_ALEVEL {
		h.log.Warnf("%s: user access level %d is not admin", loc, usrAccessLevel)
		response.ErrYouArntAdmin(w)
		return
	}

	sid, id, err := lib.GetID(api.LETTER_POINT, r.RequestURI)
	if err != nil {
		h.log.Errorf("%s: invalid id '%s' in URI '%s': %v", loc, sid, r.RequestURI, err)
		response.ErrInvalidID(w, sid)
		return
	}

	resp, err := h.crypto.Decrypt(r)
	if err != nil {
		h.log.Errorf("%s: error cannot decrypt data: %v", loc, err)
		response.ErrCannotDecryptData(w)
		return
	}
	var letter models.Letter
	if err := lib.FetchLetterFromB64(&letter, resp.Content); err != nil {
		h.log.Errorf("%s: error cannot fetch letter from b64: %v", loc, err)
		response.ErrCannotFetchFromB64(w)
		return
	}

	letter.ID = id
	h.log.Debugf("%s: decoded letter for update: %+v", loc, letter)

	msg, ok := lib.ValidateLetter(letter, h.log)
	if !ok {
		response.ErrInvalidRequest(w, msg)
		return
	}

	h.log.Debugf("%s: updating letter with id: %d", loc, id)
	if err := h.st.UpdateLetter(letter); err != nil {
		if err.Error() == st.ErrLetterNotFound(id).Error() {
			h.log.Warnf("%s: letter not found with id: %d", loc, id)
			response.ErrNotFound(w, sid)
		} else {
			h.log.Errorf("%s: failed to update letter with id: %d, error: %v", loc, id, err)
			response.ErrCannotUpdate(w)
		}
		return
	}

	h.log.Debugf("%s: letter with id %d updated successfully", loc, id)
	response.OkWithContent(w, nil)
}

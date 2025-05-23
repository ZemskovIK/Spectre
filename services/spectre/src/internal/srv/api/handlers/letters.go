package handlers

// ! TODO : think about copy-paste

import (
	"encoding/json"
	"net/http"
	"spectre/internal/lib"
	"spectre/internal/models"
	"spectre/internal/srv/api"
	"spectre/internal/srv/response"
	st "spectre/internal/storage"
	"spectre/pkg/logger"
)

const GLOC_LTS = "src/internal/api/handlers/letters.go"

type lettersStore interface {
	GetLetterByID(id int) (models.Letter, error)
	SaveLetter(letter models.Letter) error
	DeleteLetter(id int) error
	UpdateLetter(letter models.Letter) error
	GetAllLettersWithAccess(accessLevel int) ([]models.Letter, error)
}

type lettersHandler struct {
	st  lettersStore
	log *logger.Logger
}

func NewLettersHandler(
	s lettersStore, log *logger.Logger,
) *lettersHandler {
	return &lettersHandler{
		st:  s,
		log: log,
	}
}

func (h *lettersHandler) GetAll(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_LTS + "getAll()"
	h.log.Infof("%s: retrieving all letters", loc)

	accessLevel, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		h.log.Errorf("%s: access level not found or wrong type", loc)
		response.ErrFailedToRetrieveLetters(w)
		return
	}

	letters, err := h.st.GetAllLettersWithAccess(accessLevel)
	if err != nil {
		h.log.Errorf("%s: failed to retrieve letters: %v", loc, err)
		response.ErrFailedToRetrieveLetters(w)
		return
	}

	if len(letters) == 0 {
		h.log.Warnf("%s: no letters found", loc)
		response.Ok(w, []interface{}{})
		return
	}

	response.Ok(w, letters)
}

func (h *lettersHandler) GetOne(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_LTS + "getOne()"
	h.log.Infof("%s: retrieving a single letter", loc)

	usrAccess, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		response.ErrFailedToRetrieveLetters(w)
		return
	}

	sid, id, err := lib.GetID(api.LETTER_POINT, r.RequestURI)
	if err != nil {
		h.log.Errorf("%s: invalid id %s : %v", loc, sid, err)
		response.ErrInvalidID(w, sid)
		return
	}

	h.log.Infof("%s: retrieving letter with id: %d", loc, id)
	letter, err := h.st.GetLetterByID(id)
	if err != nil {
		if err.Error() == st.ErrLetterNotFound(id).Error() {
			h.log.Warnf("%s: letter not found with id: %d", loc, id)
			response.ErrNotFound(w, sid)
		} else {
			h.log.Errorf("%s: failed to retrieve letter with id: %d, error: %v",
				loc, id, err)
			response.ErrCannotGetWithID(w, sid)
		}
		return
	}

	if usrAccess < letter.AccessLevel {
		response.ErrBlockedToGet(w, usrAccess, letter.AccessLevel)
		return
	}

	response.Ok(w, letter)
}

func (h *lettersHandler) Delete(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_LTS + "delete()"

	sid, id, err := lib.GetID(api.LETTER_POINT, r.RequestURI)
	if err != nil {
		h.log.Errorf("%s: invalid id %s : %v", loc, sid, err)
		response.ErrInvalidID(w, sid)
		return
	}

	if err := h.st.DeleteLetter(id); err != nil {
		if err.Error() == st.ErrLetterNotFound(id).Error() {
			h.log.Warnf("%s: letter not found with id: %d", loc, id)
			response.ErrNotFound(w, sid)
		} else {
			response.ErrCannotDeleteWithID(w, sid)
		}
		return
	}

	response.Ok(w, nil)
}

func (h *lettersHandler) Add(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC_LTS + "add()"
	h.log.Infof("%s: adding new letter", loc)

	usrAccessLevel, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		response.ErrCannotUpdate(w)
		return
	}
	if usrAccessLevel < lib.ADMIN_ALEVEL {
		response.ErrYouArntAdmin(w)
		return
	}

	var letter models.Letter
	if err := json.NewDecoder(r.Body).Decode(&letter); err != nil {
		h.log.Errorf("%s: failed to decode JSON: %v", loc, err)
		response.ErrInvalidRequest(w, "invalid JSON")
		return
	}

	// ! TODO validation func
	if letter.Body == "" {
		h.log.Warnf("%s: body cannot be empty", loc)
		response.ErrInvalidRequest(w, "body cannot be empty")
		return
	}
	if letter.Author == "" {
		h.log.Infof("%s: author is empty, setting to 'unknown'", loc)
		letter.Author = "unknown"
	}
	if letter.FoundIn == "" {
		h.log.Infof("%s: foundIn is empty, setting to 'unknown'", loc)
		letter.FoundIn = "unknown"
	}

	h.log.Infof("%s: saving letter", loc)
	if err := h.st.SaveLetter(letter); err != nil {
		h.log.Errorf("%s: failed to save letter: %v", loc, err)
		response.ErrCannotSave(w)
		return
	}

	response.Ok(w, nil)
}

func (h *lettersHandler) Update(
	w http.ResponseWriter, r *http.Request,
) {
	log := GLOC_LTS + "update()"
	h.log.Infof("%s: updating letter", log)

	usrAccessLevel, ok := lib.FetchAccessLevelFromCtx(r.Context())
	if !ok {
		response.ErrCannotUpdate(w)
		return
	}
	if usrAccessLevel < lib.ADMIN_ALEVEL {
		response.ErrYouArntAdmin(w)
		return
	}

	sid, id, err := lib.GetID(api.LETTER_POINT, r.RequestURI)
	if err != nil {
		h.log.Errorf("%s: invalid id %s : %v", log, sid, err)
		response.ErrInvalidID(w, sid)
		return
	}

	var letter models.Letter
	if err := json.NewDecoder(r.Body).Decode(&letter); err != nil {
		h.log.Errorf("%s: failed to decode JSON: %v", log, err)
		response.ErrInvalidRequest(w, "invalid JSON")
		return
	}
	letter.ID = id

	// ! TODO validation func
	if letter.Body == "" {
		h.log.Warnf("%s: body cannot be empty", log)
		response.ErrInvalidRequest(w, "body cannot be empty")
		return
	}
	if letter.Author == "" {
		h.log.Infof("%s: author is empty, setting to 'unknown'", log)
		letter.Author = "unknown"
	}
	if letter.FoundIn == "" {
		h.log.Infof("%s: foundIn is empty, setting to 'unknown'", log)
		letter.FoundIn = "unknown"
	}

	h.log.Infof("%s: updating letter with id: %d", log, id)
	if err := h.st.UpdateLetter(letter); err != nil {
		if err.Error() == st.ErrLetterNotFound(id).Error() {
			h.log.Warnf("%s: letter not found with id: %d", log, id)
			response.ErrNotFound(w, sid)
		} else {
			h.log.Errorf("%s: failed to update letter with id: %d, error: %v", log, id, err)
			response.ErrCannotUpdate(w)
		}
		return
	}

	response.Ok(w, nil)
}

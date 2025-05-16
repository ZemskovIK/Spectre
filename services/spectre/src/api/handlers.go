package api

import (
	"encoding/json"
	"net/http"
	"spectre/api/response"
	"spectre/internal/lib"
	st "spectre/internal/storage"
	"spectre/pkg/logger"
)

var GLOC = "src/internal/api/handlers.go/"

type lettersHandler struct {
	st  st.LettersStorage
	log *logger.Logger
}

func (h *lettersHandler) getAll(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC + "getAll()"
	h.log.Infof("%s: retrieving all letters", loc)

	letters, err := h.st.GetAll()
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

	h.log.Infof("%s: successfully retrieved letters", loc)
	response.Ok(w, letters)
}

func (h *lettersHandler) getOne(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC + "getOne()"
	h.log.Infof("%s: retrieving a single letter", loc)

	sid, id, err := lib.GetID(LETTER_POINT, r.RequestURI)
	if err != nil {
		h.log.Errorf("%s: invalid id %s : %v", loc, sid, err)
		response.ErrInvalidID(w, sid)
		return
	}

	h.log.Infof("%s: retrieving letter with id: %d", loc, id)
	letter, err := h.st.Get(id)
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

	h.log.Infof("%s: successfully retrieved letter with id: %d", loc, id)
	response.Ok(w, letter)
}

func (h *lettersHandler) delete(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC + "delete()"

	sid, id, err := lib.GetID(LETTER_POINT, r.RequestURI)
	if err != nil {
		h.log.Errorf("%s: invalid id %s : %v", loc, sid, err)
		response.ErrInvalidID(w, sid)
		return
	}

	if err := h.st.Delete(id); err != nil {
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

func (h *lettersHandler) add(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC + "add()"
	_ = loc

	var letter st.Letter
	if err := json.NewDecoder(r.Body).Decode(&letter); err != nil {
		response.ErrInvalidRequest(w, "invalid JSON")
		return
	}

	// ! TODO validation func
	if letter.Body == "" {
		response.ErrInvalidRequest(w, "body cannot be empty")
		return
	}
	if letter.Author == "" {
		letter.Author = "unknown"
	}
	if letter.FoundIn == "" {
		letter.FoundIn = "unknown"
	}

	if err := h.st.Save(letter); err != nil {
		response.ErrCannotSave(w)
		return
	}

	response.Ok(w, nil)
}

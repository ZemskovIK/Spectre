package api

import (
	"encoding/json"
	"net/http"
	"spectre/api/response"
	"spectre/internal/lib"
	"spectre/internal/storage"
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
	h.log.Debugf("%s: handling getAll request", loc)
	h.log.Infof("%s: retrieving all letters", loc)

	letters, err := h.st.GetAll()
	if err != nil {
		h.log.Errorf("%s: failed to retrieve letters: %v", loc, err)
		response.ErrFailedToRetrieveLetters(w)
		return
	}

	h.log.Debugf("%s: letters retrieved: count=%d", loc, len(letters))
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
	h.log.Debugf("%s: handling getOne request, URI: %s", loc, r.RequestURI)
	h.log.Infof("%s: retrieving a single letter", loc)

	sid, id, err := lib.GetID(LETTER_POINT, r.RequestURI)
	h.log.Debugf("%s: parsed id: sid=%s, id=%d, err=%v", loc, sid, id, err)
	if err != nil {
		h.log.Errorf("%s: invalid id %s : %v", loc, sid, err)
		response.ErrInvalidID(w, sid)
		return
	}

	h.log.Debugf("%s: calling storage.Get with id=%d", loc, id)
	h.log.Infof("%s: retrieving letter with id: %d", loc, id)
	letter, err := h.st.Get(id)
	if err != nil {
		if err.Error() == storage.ErrLetterNotFound(id).Error() {
			h.log.Warnf("%s: letter not found with id: %d", loc, id)
			response.ErrNotFound(w, sid)
		} else {
			h.log.Errorf("%s: failed to retrieve letter with id: %d, error: %v",
				loc, id, err)
			response.ErrCannotGetWithID(w, sid)
		}
		return
	}

	h.log.Debugf("%s: letter retrieved: %+v", loc, letter)
	h.log.Infof("%s: successfully retrieved letter with id: %d", loc, id)
	response.Ok(w, letter)
}

func (h *lettersHandler) delete(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC + "delete()"
	h.log.Debugf("%s: handling delete request, URI: %s", loc, r.RequestURI)

	sid, id, err := lib.GetID(LETTER_POINT, r.RequestURI)
	h.log.Debugf("%s: parsed id: sid=%s, id=%d, err=%v", loc, sid, id, err)
	if err != nil {
		h.log.Errorf("%s: invalid id %s : %v", loc, sid, err)
		response.ErrInvalidID(w, sid)
		return
	}

	h.log.Debugf("%s: calling storage.Delete with id=%d", loc, id)
	if err := h.st.Delete(id); err != nil {
		if err.Error() == storage.ErrLetterNotFound(id).Error() {
			h.log.Warnf("%s: letter not found with id: %d", loc, id)
			response.ErrNotFound(w, sid)
		} else {
			h.log.Errorf("%s: failed to delete letter with id: %d, error: %v", loc, id, err)
			response.ErrCannotDeleteWithID(w, sid)
		}
		return
	}

	h.log.Infof("%s: successfully deleted letter with id: %d", loc, id)
	response.Ok(w, nil)
}

func (h *lettersHandler) add(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC + "add()"
	h.log.Debugf("%s: handling add request", loc)

	var letter st.Letter
	if err := json.NewDecoder(r.Body).Decode(&letter); err != nil {
		h.log.Errorf("%s: failed to decode request body: %v", loc, err)
		response.ErrInvalidRequest(w, "invalid JSON")
		return
	}
	h.log.Debugf("%s: decoded letter: %+v", loc, letter)

	// ! TODO validation func
	if letter.Body == "" {
		h.log.Warnf("%s: validation failed: body is empty", loc)
		response.ErrInvalidRequest(w, "body cannot be empty")
		return
	}
	if letter.Author == "" {
		h.log.Debugf("%s: author is empty, setting to 'unknown'", loc)
		letter.Author = "unknown"
	}
	if letter.FoundIn == "" {
		h.log.Debugf("%s: found_in is empty, setting to 'unknown'", loc)
		letter.FoundIn = "unknown"
	}

	h.log.Debugf("%s: calling storage.Save with letter: %+v", loc, letter)
	if err := h.st.Save(letter); err != nil {
		h.log.Errorf("%s: failed to save letter: %v", loc, err)
		response.ErrCannotSave(w)
		return
	}

	h.log.Infof("%s: successfully saved letter", loc)
	response.Ok(w, nil)
}

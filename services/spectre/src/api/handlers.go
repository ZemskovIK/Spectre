package api

import (
	"encoding/json"
	"net/http"
	"spectre/internal/storage"
	"spectre/pkg/logger"
	"strconv"
	"strings"
)

var GLOC = "src/internal/api/handlers.go/"

type lettersHandler struct {
	st  storage.LettersStorage
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
		errFailedToRetrieveLatters(w)
		return
	}

	if len(letters) == 0 {
		h.log.Infof("%s: no letters found", loc)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	h.log.Infof("%s: successfully retrieved letters", loc)
	json.NewEncoder(w).Encode(letters)
}

func (h *lettersHandler) getOne(
	w http.ResponseWriter, r *http.Request,
) {
	loc := GLOC + "getOne()"
	h.log.Infof("%s: retrieving a single letter", loc)

	sid := strings.TrimPrefix(r.URL.Path, "/api/letters/")
	id, err := strconv.Atoi(sid)
	if err != nil {
		h.log.Errorf("%s: invalid id: %s", loc, sid)
		errInvalidID(w, sid)
		return
	}

	h.log.Infof("%s: retrieving letter with id: %d", loc, id)
	letter, err := h.st.Get(id)
	if err != nil {
		if err.Error() == storage.ErrLetterNotFound(id).Error() {
			h.log.Warnf("%s: letter not found with id: %d", loc, id)
			errNotFound(w, id)
		} else {
			h.log.Errorf("%s: failed to retrieve letter with id: %d, error: %v",
				loc, id, err)
			errCannotGetWithID(w, id)
		}
		return
	}

	h.log.Infof("%s: successfully retrieved letter with id: %d", loc, id)
	json.NewEncoder(w).Encode(letter)
}

package api

import (
	"fmt"
	"net/http"
)

var (
	errFailedToRetrieveLatters = func(w http.ResponseWriter) {
		http.Error(
			w, `{"error": "failed to retrieve letters"}`,
			http.StatusInternalServerError,
		)
	}

	errInvalidID = func(w http.ResponseWriter, sid string) {
		http.Error(
			w, fmt.Sprintf(`{"error": "invalid id: %s!"}`, sid),
			http.StatusBadRequest,
		)
	}

	errCannotGetWithID = func(w http.ResponseWriter, id int) {
		http.Error(
			w, fmt.Sprintf(`{"error": "cannot get with id: %d"}`, id),
			http.StatusInternalServerError,
		)
	}

	errNotFound = func(w http.ResponseWriter, id int) {
		http.Error(
			w, fmt.Sprintf(`{"error": "letter with id %d not found"}`, id),
			http.StatusNotFound,
		)
	}
)

package api

import (
	"net/http"
	"spectre/api/methods"
	st "spectre/internal/storage"
	"spectre/pkg/logger"
)

const (
	LETTERS_POINT = "/api/letters"
	LETTER_POINT  = "/api/letters/"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(s st.LettersStorage, log *logger.Logger) *Router {
	mux := http.NewServeMux()

	lettersHL := lettersHandler{
		st:  s,
		log: log,
	}

	mux.Handle(methods.GET(LETTERS_POINT),
		CORSMiddleware(
			JSONRespMiddleware(
				http.HandlerFunc(lettersHL.getAll),
			),
		))
	mux.Handle(methods.GET(LETTER_POINT),
		CORSMiddleware(
			JSONRespMiddleware(
				http.HandlerFunc(lettersHL.getOne),
			),
		))

	return &Router{
		mux: mux,
	}
}

func (rt *Router) ServeHTTP(
	w http.ResponseWriter, r *http.Request,
) {
	if rt.mux == nil {
		http.Error(w, "Router is not initialized", http.StatusInternalServerError)
		return
	}
	rt.mux.ServeHTTP(w, r)
}

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
	mux http.Handler
}

func NewRouter(s st.LettersStorage, log *logger.Logger) *Router {
	mux := http.NewServeMux()

	lettersHL := lettersHandler{
		st:  s,
		log: log,
	}

	// CORS
	mux.Handle(methods.OPTIONS(LETTER_POINT),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	mux.Handle(methods.OPTIONS(LETTERS_POINT),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)

	mux.Handle(methods.GET(LETTERS_POINT),
		http.HandlerFunc(lettersHL.getAll),
	)
	mux.Handle(methods.GET(LETTER_POINT),
		http.HandlerFunc(lettersHL.getOne),
	)
	mux.Handle(methods.PUT(LETTER_POINT),
		http.HandlerFunc(lettersHL.update),
	)
	mux.Handle(methods.DELETE(LETTER_POINT),
		http.HandlerFunc(lettersHL.delete),
	)
	mux.Handle(methods.POST(LETTERS_POINT),
		http.HandlerFunc(lettersHL.add),
	)

	mwmux := CORSMiddleware(
		JSONRespMiddleware(
			mux,
		),
	)

	return &Router{
		mux: mwmux,
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

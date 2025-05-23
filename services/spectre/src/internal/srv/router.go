package server

import (
	"net/http"

	"spectre/internal/srv/api"
	"spectre/internal/srv/api/handlers"
	"spectre/internal/srv/auth"
	"spectre/internal/srv/methods"
	st "spectre/internal/storage"
	"spectre/pkg/logger"
)

const (
	LOGIN_POINT = "/login"
)

type Router struct {
	mux http.Handler
}

func NewRouter(s st.Storage, log *logger.Logger) *Router {
	mux := http.NewServeMux()

	lettersHL := handlers.NewLettersHandler(s, log)
	authHL := auth.NewAuthHandler(s, log)

	// CORS
	mux.Handle(methods.OPTIONS(api.LETTER_POINT),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	mux.Handle(methods.OPTIONS(api.LETTERS_POINT),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)

	// api
	mux.Handle(methods.GET(api.LETTERS_POINT),
		http.HandlerFunc(lettersHL.GetAll),
	)
	mux.Handle(methods.GET(api.LETTER_POINT),
		http.HandlerFunc(lettersHL.GetOne),
	)
	mux.Handle(methods.PUT(api.LETTER_POINT),
		http.HandlerFunc(lettersHL.Update),
	)
	mux.Handle(methods.DELETE(api.LETTER_POINT),
		http.HandlerFunc(lettersHL.Delete),
	)
	mux.Handle(methods.POST(api.LETTERS_POINT),
		http.HandlerFunc(lettersHL.Add),
	)

	// auth
	mux.Handle(methods.POST(auth.LOGIN_POINT),
		http.HandlerFunc(authHL.Login),
	)

	return &Router{
		mux: mux,
	}
}

func (rt *Router) Use(mw func(http.Handler) http.Handler) {
	rt.mux = mw(rt.mux)
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

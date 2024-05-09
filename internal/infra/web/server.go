package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewWebServer(serverPort string) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// TODO: add rate limiter middleware
	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
	})

	err := http.ListenAndServe(serverPort, r)
	if err != nil {
		log.Println(err)
	}
}

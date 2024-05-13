package web

import (
	"log"
	"net/http"

	"github.com/flpnascto/rate-limiter-go/configs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewWebServer(cfg *configs.Conf) {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(RateLimiterMiddleware(cfg))
	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
	})

	err := http.ListenAndServe(cfg.WebServerPort, r)
	if err != nil {
		log.Println(err)
	}
}

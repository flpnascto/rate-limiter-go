package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/flpnascto/rate-limiter-go/internal/entity"
	"github.com/flpnascto/rate-limiter-go/internal/usecases"
	"github.com/spf13/viper"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewWebServer(db entity.LimiterRepositoryInterface) {
	log.Println("Starting server...")
	r := chi.NewRouter()
	r.Use(middleware.RealIP)

	r.Use(LimiterMiddleware(db))
	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		})
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			response := usecases.NewListUseCase(db).Execute()
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				log.Println(err)
			}
		})
	})

	err := http.ListenAndServe(viper.GetString("WebServerPort"), r)
	if err != nil {
		log.Println(err)
	}
}

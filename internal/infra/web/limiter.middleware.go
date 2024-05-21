package web

import (
	"log"
	"net"
	"net/http"

	"github.com/flpnascto/rate-limiter-go/internal/entity"
	"github.com/flpnascto/rate-limiter-go/internal/usecases"
)

func LimiterMiddleware(LimiterRepository entity.LimiterRepositoryInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("LimiterMiddleware - Iniciado")
			ip := r.RemoteAddr
			if ip != "" {
				host, _, err := net.SplitHostPort(ip)
				if err != nil {
					log.Println("LimiterMiddleware - Error 1:", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				ip = host
			}
			token := r.Header.Get("API_KEY")
			log.Println("IP:", ip, "TOKEN:", token)
			limiter := entity.NewLimiter(ip, &token)
			log.Println("Limiter:", limiter)
			err := usecases.NewVerifyLimiterUseCase(LimiterRepository).Execute(limiter)
			if err != nil {
				if err.Error() == "Rate limit exceeded" {
					http.Error(w, err.Error(), http.StatusTooManyRequests)
					return
				}
				log.Println("LimiterMiddleware - Error 2:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

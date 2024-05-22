package web

import (
	"net"
	"net/http"

	"github.com/flpnascto/rate-limiter-go/internal/entity"
	"github.com/flpnascto/rate-limiter-go/internal/usecases"
)

func LimiterMiddleware(LimiterRepository entity.LimiterRepositoryInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			manyRequestMessage := "you have reached the maximum number of requests or actions allowed within a certain time frame"
			ip := r.RemoteAddr
			if ip != "" {
				host, _, err := net.SplitHostPort(ip)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				ip = host
			}
			token := r.Header.Get("API_KEY")
			limiter := entity.NewLimiter(ip, token)
			err := usecases.NewVerifyLimiterUseCase(LimiterRepository).Execute(limiter)
			if err != nil {
				if err.Error() == "Rate limit exceeded" {
					http.Error(w, manyRequestMessage, http.StatusTooManyRequests)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

package web

import (
	"log"
	"net"
	"net/http"

	"github.com/flpnascto/rate-limiter-go/internal/entity"
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	ipRateLimiter := entity.NewRateLimiter(5)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if ip != "" {
			host, _, err := net.SplitHostPort(ip)
			if err != nil {
				log.Fatal(err)
			}
			ip = host
		}
		err := ipRateLimiter.AddIpVisitor(ip)
		if err != nil {
			message := "too many requests from " + ip
			http.Error(w, message, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

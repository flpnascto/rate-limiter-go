package web

import (
	"log"
	"net"
	"net/http"

	"github.com/flpnascto/rate-limiter-go/configs"
	"github.com/flpnascto/rate-limiter-go/internal/entity"
)

func RateLimiterMiddleware(cfg *configs.Conf) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		ipRateLimiter := entity.NewRateLimiter(cfg.MaxIpRequests, cfg.IpBlockTime)
		tokenRateLimiter := entity.NewRateLimiter(cfg.MaxTokenRequests, cfg.TokenBlockTime)
		var message string
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if apiKey := r.Header.Get("API_KEY"); apiKey != "" {
				message = requestWithToken(r, tokenRateLimiter)
			} else {
				message = requestWithIp(r, ipRateLimiter)
			}
			if message != "" {
				http.Error(w, message, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func requestWithToken(r *http.Request, ipRateLimiter *entity.RateLimiter) string {
	token := r.Header.Get("API_KEY")
	err := ipRateLimiter.AddIpVisitor(token)
	if err != nil {
		//TODO: change message for: "you have reached the maximum number of requests or actions allowed within a certain time frame"
		message := "too many requests with token " + token
		return message
	}
	return ""
}

func requestWithIp(r *http.Request, tokenRateLimiter *entity.RateLimiter) string {
	ip := r.RemoteAddr
	if ip != "" {
		host, _, err := net.SplitHostPort(ip)
		if err != nil {
			log.Fatal(err)
		}
		ip = host
	}
	err := tokenRateLimiter.AddIpVisitor(ip)
	if err != nil {
		//TODO: change message for: "you have reached the maximum number of requests or actions allowed within a certain time frame"
		message := "too many requests from " + ip
		return message
	}
	return ""
}

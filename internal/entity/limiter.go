package entity

import "time"

type Limiter struct {
	Ip       string
	Token    *string
	Requests int8
	Block    bool
	Time     time.Time
}

func NewLimiter(ip string, token *string) *Limiter {
	return &Limiter{
		Ip:       ip,
		Token:    token,
		Requests: 1,
		Block:    false,
		Time:     time.Now(),
	}
}

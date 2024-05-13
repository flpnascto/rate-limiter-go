package entity

import (
	"errors"
	"sync"
	"time"
)

type Visitor struct {
	identity string
	requests int
	created  time.Time
}

type RateLimiter struct {
	visitors    []Visitor
	maxRequests int
	blockTime   int
	mu          sync.Mutex
}

func NewRateLimiter(limit int, block int) *RateLimiter {

	r := &RateLimiter{
		visitors:    []Visitor{},
		maxRequests: limit,
		blockTime:   block,
	}

	go r.cleanup()

	return r
}

func (r *RateLimiter) AddIpVisitor(ip string) error {
	index := -1
	for i, v := range r.visitors {
		if v.identity == ip {
			index = i
			break
		}
	}
	if index == -1 {
		r.visitors = append(r.visitors, newVisitor(ip))
	} else {
		r.visitors[index].requests++
		r.visitors[index].created = time.Now()
		err := r.limit(r.visitors[index])
		if err != nil {
			return err
		}
	}
	return nil
}

func newVisitor(ip string) Visitor {
	return Visitor{
		identity: ip,
		requests: 1,
		created:  time.Now(),
	}
}

func (r *RateLimiter) limit(v Visitor) error {
	if v.requests > r.maxRequests {
		return errors.New("too many requests")
	}
	return nil
}

func (r *RateLimiter) cleanup() {
	for {
		time.Sleep(100 * time.Millisecond)

		r.mu.Lock()
		newVisitors := make([]Visitor, 0)
		for _, y := range r.visitors {
			if y.requests > r.maxRequests {
				if time.Since(y.created) > time.Duration(r.blockTime)*time.Second {
					newVisitors = append(newVisitors, y)
				}
			} else {
				if time.Since(y.created) <= 1*time.Second {
					newVisitors = append(newVisitors, y)
				}
			}
		}
		r.visitors = newVisitors
		r.mu.Unlock()
	}
}

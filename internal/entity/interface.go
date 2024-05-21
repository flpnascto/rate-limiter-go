package entity

import "time"

type LimiterRepositoryInterface interface {
	Create(limiter *Limiter, duration time.Duration) error
	GetByIp(ip string) (Limiter, error)
	Update(limiter *Limiter, duration time.Duration) error
	Delete(ip string) error
	List() ([]Limiter, error)
}

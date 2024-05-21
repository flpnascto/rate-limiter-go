package usecases

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/flpnascto/rate-limiter-go/internal/entity"
	"github.com/spf13/viper"
)

type VerifyLimiterUseCase struct {
	LimiterRepository entity.LimiterRepositoryInterface
	mu                sync.Mutex
}

func NewVerifyLimiterUseCase(limiterRepository entity.LimiterRepositoryInterface) *VerifyLimiterUseCase {
	return &VerifyLimiterUseCase{
		LimiterRepository: limiterRepository,
	}
}

func (c *VerifyLimiterUseCase) Execute(limiter *entity.Limiter) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	result, _ := c.LimiterRepository.GetByIp(limiter.Ip)

	if c.limitExceeded(&result) {
		result.Block = true
		c.LimiterRepository.Update(&result, c.setDuration(limiter))
		return errors.New("Rate limit exceeded")
	}

	if result.Ip == "" {
		err := c.LimiterRepository.Create(limiter, c.setDuration(limiter))
		if err != nil {
			log.Println("UseCase error 2:", err)
			return err
		}
		return nil
	}
	result.Requests++
	err := c.LimiterRepository.Update(&result, c.setDuration(limiter))
	if err != nil {
		return err
	}
	return nil
}

func (c *VerifyLimiterUseCase) limitExceeded(limiter *entity.Limiter) bool {
	if limiter.Block {
		return true
	}
	var maxRequests int
	if limiter.Token == "" {
		maxRequests = viper.GetInt("MaxIpRequests")
	} else {
		maxRequests = viper.GetInt("MaxTokenRequests")
	}
	return limiter.Requests > int8(maxRequests)
}

func (c *VerifyLimiterUseCase) setDuration(limiter *entity.Limiter) time.Duration {
	if !limiter.Block {
		return time.Second
	}
	var quantity int
	if limiter.Token == "" {
		quantity = viper.GetInt("TokenBlockDuration")
	} else {
		quantity = viper.GetInt("IpBlockDuration")
	}
	return time.Duration(quantity) * time.Second
}

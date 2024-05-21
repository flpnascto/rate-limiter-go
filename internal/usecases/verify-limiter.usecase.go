package usecases

import (
	"errors"
	"log"
	"time"

	"github.com/flpnascto/rate-limiter-go/internal/entity"
	"github.com/spf13/viper"
)

type VerifyLimiterUseCase struct {
	LimiterRepository entity.LimiterRepositoryInterface
}

func NewVerifyLimiterUseCase(limiterRepository entity.LimiterRepositoryInterface) *VerifyLimiterUseCase {
	return &VerifyLimiterUseCase{
		LimiterRepository: limiterRepository,
	}
}

func (c *VerifyLimiterUseCase) Execute(limiter *entity.Limiter) error {
	log.Println("UseCase:", limiter)
	var err error
	result, err := c.LimiterRepository.GetByIp(limiter.Ip)
	log.Println("UseCase result:", result)
	if err != nil {
		log.Println("UseCase error 1:", err)
	}

	if result.Ip == "" {
		err = c.LimiterRepository.Create(limiter, c.setDuration(limiter))
		if err != nil {
			log.Println("UseCase error 2:", err)
			return err
		}
		return nil
	} else {
		result.Requests++
		if c.limitExceeded(&result) {
			result.Block = true
			log.Println("UseCase error exceeded :", result)
			return errors.New("too many requests")
		}
		log.Println("UseCase result:", &result)
		err = c.LimiterRepository.Update(&result, c.setDuration(limiter))
		log.Println("UseCase error 3:", err)
		return err
	}
	return nil
}

func (c *VerifyLimiterUseCase) limitExceeded(limiter *entity.Limiter) bool {
	if limiter.Block {
		return true
	}
	var maxRequests int
	if limiter.Token != nil {
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
	if limiter.Token != nil {
		quantity = viper.GetInt("TokenBlockDuration")
	} else {
		quantity = viper.GetInt("IpBlockDuration")
	}
	return time.Duration(quantity) * time.Second
}

package usecases

import (
	"github.com/flpnascto/rate-limiter-go/internal/entity"
)

type ListUseCase struct {
	LimiterRepository entity.LimiterRepositoryInterface
}

func NewListUseCase(limiterRepository entity.LimiterRepositoryInterface) *ListUseCase {
	return &ListUseCase{
		LimiterRepository: limiterRepository,
	}
}

func (c *ListUseCase) Execute() []entity.Limiter {
	var err error
	result, err := c.LimiterRepository.List()
	if err != nil {
		return nil
	}
	return result
}

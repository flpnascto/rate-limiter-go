package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/flpnascto/rate-limiter-go/internal/entity"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Db *redis.Client
}

func NewRedisClient(db *redis.Client) *RedisClient {
	return &RedisClient{Db: db}
}

var ctx = context.Background()

func (r *RedisClient) Create(limiter *entity.Limiter, duration time.Duration) error {
	log.Println("REDIS CREATE -> limiter:", limiter)
	data, err := json.Marshal(limiter)
	if err != nil {
		return fmt.Errorf("failed to serialize limiter: %w", err)
	}

	err = r.Db.Set(ctx, limiter.Ip, string(data), duration).Err()
	if err != nil {
		return fmt.Errorf("failed to save limiter to Redis: %w", err)
	}

	return nil
}

func (r *RedisClient) GetByIp(ip string) (entity.Limiter, error) {
	log.Println("REDIS GET BY IP -> ip:", ip)
	val, err := r.Db.Get(ctx, ip).Result()
	log.Println("REDIS GET BY IP -> val:", val)

	if err == redis.Nil {
		return entity.Limiter{}, fmt.Errorf("no limiter found for IP: %s", ip)
	} else if err != nil {
		return entity.Limiter{}, fmt.Errorf("failed to get limiter from Redis: %w", err)
	}

	var limiter entity.Limiter
	err = json.Unmarshal([]byte(val), &limiter)
	if err != nil {
		return entity.Limiter{}, fmt.Errorf("failed to deserialize limiter: %w", err)
	}
	log.Println("REDIS GET BY IP -> limiter:", limiter)

	return limiter, nil
}

func (r *RedisClient) Update(limiter *entity.Limiter, duration time.Duration) error {
	log.Println("REDIS UPDATE -> limiter:", limiter)
	data, err := json.Marshal(limiter)
	if err != nil {
		return fmt.Errorf("failed to serialize limiter: %w", err)
	}

	err = r.Db.Set(ctx, limiter.Ip, string(data), duration).Err()
	if err != nil {
		return fmt.Errorf("failed to save limiter to Redis: %w", err)
	}

	return nil
}

func (r *RedisClient) Delete(ip string) error {
	return nil
}

func (r *RedisClient) List() ([]entity.Limiter, error) {
	keys, err := r.Db.Keys(ctx, "*").Result()
	log.Println("REDIS LIST -> keys:", keys)
	if err != nil {
		return nil, fmt.Errorf("failed to get keys from Redis: %w", err)
	}

	results := make(map[string]string)
	for _, key := range keys {
		value, err := r.Db.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				continue // key does not exist
			} else {
				return nil, fmt.Errorf("failed to get value for key %s from Redis: %w", key, err)
			}
		}
		results[key] = value
	}
	var limiterList []entity.Limiter

	for _, item := range results {
		var limiter entity.Limiter
		err := json.Unmarshal([]byte(item), &limiter)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize limiter: %w", err)
		}
		limiterList = append(limiterList, limiter)
	}
	return limiterList, nil
}

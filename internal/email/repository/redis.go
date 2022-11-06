package repository

import "github.com/go-redis/redis"

type EmailRepository struct {
	redis *redis.Client
}

func NewEmailRepository(redis *redis.Client) *EmailRepository {
	return &EmailRepository{
		redis: redis,
	}
}

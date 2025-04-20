package di

import "time"

type IRedisResetRepository interface { // not used on this project(redirected to notifications)
	SaveToken(email, code string, expiresAt time.Time) error
	GetToken(email string) (string, time.Time, error)
	DeleteToken(email string) error
	GetResetCodeCount(email string) (int, error)
	IncrementResetCodeCount(email string, ttl time.Duration) error
}

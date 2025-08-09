package middleware

import (
	"purpleschool/internal/auth"
	"purpleschool/internal/logger"
)

type Middleware struct {
	logger *logger.Logger
	jwt    *auth.JWT
}

func NewMiddleware(l *logger.Logger, jwt *auth.JWT) *Middleware {
	return &Middleware{
		logger: l,
		jwt:    jwt,
	}
}

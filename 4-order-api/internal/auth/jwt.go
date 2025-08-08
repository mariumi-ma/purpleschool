package auth

import "github.com/golang-jwt/jwt/v5"

type JWT struct {
	SecretKey string
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		SecretKey: secretKey,
	}
}

func (j *JWT) GenerateToken(sessionID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_id": sessionID,
	})

	return token.SignedString([]byte(j.SecretKey))
}

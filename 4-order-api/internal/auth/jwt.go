package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	SecretKey string
}

type JWTData struct {
	UserID uint
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		SecretKey: secretKey,
	}
}

func (j *JWT) GenerateToken(data JWTData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": data.UserID,
	})

	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWT) ParseToken(token string) (bool, *JWTData) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return false, nil
	}

	// Обработка userID из токена
	userID, ok := tokenClaims.Claims.(jwt.MapClaims)["user_id"]
	if !ok {
		return false, nil
	}

	userIDFloat, ok := userID.(float64)
	if !ok {
		return false, nil
	}

	return tokenClaims.Valid, &JWTData{
		UserID: uint(userIDFloat),
	}
}

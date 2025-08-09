package auth

import "github.com/golang-jwt/jwt/v5"

type JWT struct {
	SecretKey string
}

type JWTData struct {
	Phone string
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		SecretKey: secretKey,
	}
}

func (j *JWT) GenerateToken(data JWTData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": data.Phone,
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

	phone := tokenClaims.Claims.(jwt.MapClaims)["phone"].(string)

	return tokenClaims.Valid, &JWTData{
		Phone: phone,
	}
}

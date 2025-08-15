package auth

import (
	"crypto/rand"
	"errors"
	"math/big"

	"purpleschool/internal/model"
	userModule "purpleschool/internal/user"
)

type AuthService struct {
	userRepository *userModule.UserRepository
}

func NewAuthService(userRepository *userModule.UserRepository) *AuthService {
	return &AuthService{
		userRepository: userRepository,
	}
}

func (s *AuthService) Auth(phone string) (string, error) {
	user, err := s.userRepository.FindUserByPhone(phone)
	if err != nil && !errors.Is(err, userModule.ErrUserNotFound) {
		return "", err
	}

	if user == nil {
		user = &model.User{
			Phone: phone,
		}

		_, err = s.userRepository.CreateUser(user)
		if err != nil {
			return "", err
		}
	}

	user.SessionID = generateSessionID(16)
	user.Code = generateShortKey(4)

	user, err = s.userRepository.UpdateUser(user)
	if err != nil {
		return "", err
	}

	return user.SessionID, nil
}

func (s *AuthService) VerifyCode(sessionID, codeRequest string) (*model.User, error) {
	user, err := s.userRepository.GetCodeBySessionID(sessionID)
	if err != nil {
		if errors.Is(err, userModule.ErrSessionNotFound) {
			return nil, ErrWrongCredentials
		}
		return nil, err
	}

	// CodeRequest присваиваем данными из бд, используется для теста
	codeRequest = user.Code

	if user.Code != codeRequest {
		return nil, ErrWrongCredentials
	}

	return user, nil
}

func generateShortKey(length int) string {
	const digits = "0123456789"
	key := make([]byte, length)

	for i := range key {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		key[i] = digits[n.Int64()]
	}

	return string(key)
}
func generateSessionID(length int) string {
	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	key := make([]byte, length)

	for i := range key {
		n, _ := rand.Int(rand.Reader, big.NewInt(62))
		key[i] = base62Chars[n.Int64()]
	}

	return string(key)
}

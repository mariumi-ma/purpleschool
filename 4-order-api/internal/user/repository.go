package user

import (
	"errors"

	"purpleschool/internal/database"

	"gorm.io/gorm"
)

var ErrUserNotFound error = errors.New("user not found")
var ErrSessionNotFound error = errors.New("session not found")

type UserRepository struct {
	Database *database.Db
}

func NewUserRepository(db *database.Db) *UserRepository {
	return &UserRepository{
		Database: db,
	}
}

func (r *UserRepository) FindUserByPhone(phone string) (*User, error) {
	var user *User

	result := r.Database.First(&user, "phone = ?", phone)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}

	return user, nil
}

func (r *UserRepository) GetCodeBySessionID(sessionID string) (*User, error) {
	var user *User

	result := r.Database.First(&user, "session_id = ?", sessionID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, result.Error
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user *User) (*User, error) {
	result := r.Database.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user *User) (*User, error) {
	result := r.Database.Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

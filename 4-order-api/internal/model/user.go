package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Phone     string `gorm:"index"`
	SessionID string
	Code      string
	Orders    []Order
}

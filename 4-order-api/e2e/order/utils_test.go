package order_e2e

import (
	"fmt"
	"os"

	"testing"

	"purpleschool/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func createTestData(db *gorm.DB, t *testing.T) {
	t.Helper()

	err := db.Create(&model.User{
		Phone: "1234567890",
	}).Error
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	err = db.Create(&model.Product{
		Name:        "Product 1",
		Description: "Description 1",
	}).Error
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}
}

func getTestData(db *gorm.DB, t *testing.T) (*model.User, *model.Product) {
	t.Helper()

	var user *model.User
	if err := db.Where("phone = ?", "1234567890").First(&user).Error; err != nil {
		t.Fatalf("Failed to get test user: %v", err)
	}

	var product *model.Product
	if err := db.Where("name = ?", "Product 1").First(&product).Error; err != nil {
		t.Fatalf("Failed to get test product: %v", err)
	}

	return user, product
}

func removeData(db *gorm.DB, userID uint) {
	db.Unscoped().
		Where("phone = ?", "1234567890").
		Delete(&model.User{})

	db.Unscoped().
		Where("name = ?", "Product 1").
		Delete(&model.Product{})

	db.Unscoped().
		Where("user_id = ?", userID).
		Delete(&model.Order{})

	db.Exec(`
		DELETE FROM order_products 
		WHERE order_id IN (SELECT id FROM orders WHERE user_id = ?)
	`, userID)
}

func generateTestToken(userID uint) (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}

	secretKey := os.Getenv("TOKEN_SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})

	return token.SignedString([]byte(secretKey))
}

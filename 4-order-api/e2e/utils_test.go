package order_e2e

import (
	"fmt"
	"os"

	"testing"

	"purpleschool/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB(t *testing.T) *gorm.DB {
	if err := godotenv.Load(); err != nil {
		require.NoError(t, err, "Failed to load environment variables")
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
		require.NoError(t, err, "Failed to connect to database")
	}

	return db
}

func createTestData(db *gorm.DB, t *testing.T, createOrder bool) (*model.User, *model.Product, *model.Order) {
	t.Helper()

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	user := &model.User{
		Phone: "1234567890",
	}

	err := tx.Create(user).Error
	require.NoError(t, err, "Failed to create test user")

	product := &model.Product{
		Name:        "Product 1",
		Description: "Description 1",
	}

	err = tx.Create(product).Error
	require.NoError(t, err, "Failed to create test product")

	var order *model.Order
	if createOrder {
		order = &model.Order{
			UserID:   user.ID,
			Products: []model.Product{*product},
		}

		err := tx.Create(order).Error
		require.NoError(t, err, "Failed to create test order")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		require.NoError(t, err, "Failed to commit transaction")
	}

	return user, product, order
}

func removeData(db *gorm.DB, userID uint) {
	db.Exec(`
		DELETE FROM order_products
		WHERE order_id IN (SELECT id FROM orders WHERE user_id = ?)
	`, userID)

	db.Unscoped().Where("1 = 1").Delete(&model.Order{})
	db.Unscoped().Where("1 = 1").Delete(&model.User{})
	db.Unscoped().Where("1 = 1").Delete(&model.Product{})
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

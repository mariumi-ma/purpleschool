package order_e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"purpleschool/internal/app"
	"purpleschool/internal/ctxutils"
	"purpleschool/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	//Prepare
	db := initDB(t)
	user, product, _ := createTestData(db, t, false)
	t.Cleanup(func() {
		removeData(db, user.ID)
	})

	token, err := generateTestToken(user.ID)
	require.NoError(t, err, "Failed to generate token")

	ctx := ctxutils.WithUserID(context.Background(), user.ID)

	testCases := []struct {
		name           string
		body           model.CreateOrderRequest
		expectedStatus int
	}{
		{
			name: "Success",
			body: model.CreateOrderRequest{
				ProductIDs: []uint{product.ID},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Empty body",
			body: model.CreateOrderRequest{
				ProductIDs: []uint{},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, err := json.Marshal(tt.body)
			require.NoError(t, err, "Failed to marshal request")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/order", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			req = req.WithContext(ctx)

			ap := app.App()
			ap.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var resp model.CreateOrderResponse
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err, "Failed to unmarshal request")

				// Проверяем что заказ действительно создался в БД
				order := &model.Order{}
				err = db.First(order, "id = ?", resp.OrderID).Error
				require.NoError(t, err, "Failed to find order in DB")

				assert.Equal(t, user.ID, order.UserID)
			}
		})
	}
}

func TestGetOrderByID(t *testing.T) {
	//Prepare
	db := initDB(t)
	user, _, order := createTestData(db, t, true)
	t.Cleanup(func() {
		removeData(db, user.ID)
	})

	token, err := generateTestToken(user.ID)
	require.NoError(t, err, "Failed to generate token")

	ctx := ctxutils.WithUserID(context.Background(), user.ID)

	testCases := []struct {
		name                    string
		id                      string
		expectedStatus          int
		expectedSuccessResponse model.OrderResponse
		expectedErrorResponse   string
	}{
		{
			name:                    "Success",
			id:                      fmt.Sprintf("%d", order.ID),
			expectedStatus:          http.StatusOK,
			expectedSuccessResponse: order.ToResponse(),
		},
		{
			name:                  "Invalid ID",
			id:                    "invalid",
			expectedStatus:        http.StatusBadRequest,
			expectedErrorResponse: "invalid id",
		},
		{
			name:                  "Order Not Found",
			id:                    fmt.Sprintf("%d", 11111111),
			expectedStatus:        http.StatusNotFound,
			expectedErrorResponse: "order not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			url := fmt.Sprintf("/order/%s", tt.id)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			req = req.WithContext(ctx)

			ap := app.App()
			ap.ServeHTTP(w, req)

			var resp []byte
			if tt.expectedStatus == http.StatusOK {
				resp, err = json.Marshal(tt.expectedSuccessResponse)
				require.NoError(t, err, "Failed to unmarshal response")
			} else {
				resp, err = json.Marshal(tt.expectedErrorResponse)
				require.NoError(t, err, "Failed to unmarshal response")
			}

			require.Equal(t, tt.expectedStatus, w.Code)
			require.JSONEq(t, string(resp), w.Body.String())
		},
		)
	}
}

func TestGetOrderByUserID(t *testing.T) {
	//Prepare
	db := initDB(t)
	user, _, order := createTestData(db, t, true)
	t.Cleanup(func() {
		removeData(db, user.ID)
	})

	orders := model.Orders{}
	orders = append(orders, *order)

	token, err := generateTestToken(user.ID)
	require.NoError(t, err, "Failed to generate token")

	ctx := ctxutils.WithUserID(context.Background(), user.ID)

	testCases := []struct {
		name                    string
		expectedStatus          int
		expectedSuccessResponse model.OrdersResponse
		expectedErrorResponse   string
	}{
		{
			name:                    "Success",
			expectedStatus:          http.StatusOK,
			expectedSuccessResponse: orders.ToResponse(),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/my-orders", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			req = req.WithContext(ctx)

			ap := app.App()
			ap.ServeHTTP(w, req)

			var resp []byte
			if tt.expectedStatus == http.StatusOK {
				resp, err = json.Marshal(tt.expectedSuccessResponse)
				require.NoError(t, err, "Failed to unmarshal response")
			} else {
				resp, err = json.Marshal(tt.expectedErrorResponse)
				require.NoError(t, err, "Failed to unmarshal response")
			}

			require.Equal(t, tt.expectedStatus, w.Code)
			require.JSONEq(t, string(resp), w.Body.String())
		},
		)
	}
}

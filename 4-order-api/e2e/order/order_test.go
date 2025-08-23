package order_e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"purpleschool/internal/app"
	"purpleschool/internal/ctxutils"
	"purpleschool/internal/model"
)

func TestCreateOrderSuccess(t *testing.T) {
	//Prepare
	db := initDB()
	createTestData(db, t)

	user, product := getTestData(db, t)
	defer removeData(db, user.ID)

	dataReq, err := json.Marshal(model.CreateOrderRequest{
		ProductIDs: []uint{product.ID},
	})
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	token, err := generateTestToken(user.ID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	ctx := ctxutils.WithUserID(context.Background(), user.ID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/order", bytes.NewReader(dataReq))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req = req.WithContext(ctx)

	ap := app.App()
	ap.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected: %d, got: %d", http.StatusCreated, w.Code)
	}

	var resp model.CreateOrderResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	order := &model.Order{}
	err = db.First(order, "id = ?", resp.OrderID).Error
	if err != nil {
		t.Fatalf("Failed to find order in DB: %v", err)
	}

	if order.UserID != user.ID {
		t.Errorf("expected user id: %d, got: %d", user.ID, order.UserID)
	}
}

func TestCreateOrderEmptyBody(t *testing.T) {
	//Prepare
	db := initDB()
	createTestData(db, t)

	user, _ := getTestData(db, t)
	defer removeData(db, user.ID)

	dataReq, err := json.Marshal(model.CreateOrderRequest{
		ProductIDs: []uint{},
	})
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	token, err := generateTestToken(user.ID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	ctx := ctxutils.WithUserID(context.Background(), user.ID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/order", bytes.NewReader(dataReq))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req = req.WithContext(ctx)

	ap := app.App()
	ap.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected: %d, got: %d", http.StatusBadRequest, w.Code)
	}
}

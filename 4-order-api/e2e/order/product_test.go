package order_e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"

	"purpleschool/internal/app"
	"purpleschool/internal/model"

	"github.com/stretchr/testify/require"
)

func TestGetProductByID(t *testing.T) {
	//Prepare
	db := initDB(t)
	// TODO возвращать слайс продуктов
	user, product, _ := createTestData(db, t, true)
	t.Cleanup(func() {
		removeData(db, user.ID)
	})

	testCases := []struct {
		name                    string
		id                      string
		expectedStatus          int
		expectedSuccessResponse *model.Product
		expectedErrorResponse   string
	}{
		{
			name:                    "Success",
			id:                      fmt.Sprintf("%d", product.ID),
			expectedStatus:          http.StatusOK,
			expectedSuccessResponse: product,
		},
		{
			name:                  "Invalid ID",
			id:                    "invalid",
			expectedStatus:        http.StatusBadRequest,
			expectedErrorResponse: "invalid id",
		},
		{
			name:                  "Product Not Found",
			id:                    fmt.Sprintf("%d", 11111111),
			expectedStatus:        http.StatusNotFound,
			expectedErrorResponse: "product not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			url := fmt.Sprintf("/products/%s", tt.id)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("Content-Type", "application/json")

			ap := app.App()
			ap.ServeHTTP(w, req)

			var (
				resp []byte
				err  error
			)

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

func TestGetProducts(t *testing.T) {
	//Prepare
	db := initDB(t)
	user, product, _ := createTestData(db, t, true)
	t.Cleanup(func() {
		removeData(db, user.ID)
	})

	products := model.Products{}
	products = append(products, *product)

	testCases := []struct {
		name                    string
		expectedStatus          int
		expectedSuccessResponse model.ProductsResponse
		expectedErrorResponse   string
	}{
		{
			name:                    "Success",
			expectedStatus:          http.StatusOK,
			expectedSuccessResponse: products.ToResponse(),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/products", nil)
			req.Header.Set("Content-Type", "application/json")

			ap := app.App()
			ap.ServeHTTP(w, req)

			var (
				resp []byte
				err  error
			)

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

func TestCreateProduct(t *testing.T) {
	//Prepare
	db := initDB(t)
	user, _, _ := createTestData(db, t, false)
	t.Cleanup(func() {
		removeData(db, user.ID)
	})

	testCases := []struct {
		name           string
		body           model.CreateProductRequest
		expectedStatus int
	}{
		{
			name: "Success",
			body: model.CreateProductRequest{
				Name:        "Test Product",
				Description: "Test Description",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Empty body",
			body:           model.CreateProductRequest{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, err := json.Marshal(tt.body)
			require.NoError(t, err, "Failed to marshal request")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/products", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")

			ap := app.App()
			ap.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var resp model.ProductResponse
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err, "Failed to unmarshal request")

				// Проверяем что товар действительно создался в БД
				product := &model.Product{}
				err = db.First(product, "id = ?", resp.ID).Error
				require.NoError(t, err, "Failed to find product in DB")
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	//Prepare
	db := initDB(t)
	user, product, _ := createTestData(db, t, false)
	t.Cleanup(func() {
		removeData(db, user.ID)
	})

	testCases := []struct {
		name                  string
		id                    string
		body                  model.UpdateProductRequest
		expectedStatus        int
		expectedErrorResponse string
	}{
		{
			name: "Success",
			id:   fmt.Sprintf("%d", product.ID),
			body: model.UpdateProductRequest{
				Name:        "Test Updated Product",
				Description: "Test Updated Description",
			},
			expectedStatus: http.StatusOK,
		},
		{
			// TODO check
			name:           "Empty body",
			id:             fmt.Sprintf("%d", product.ID),
			body:           model.UpdateProductRequest{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:                  "invalid id",
			id:                    "invalid",
			body:                  model.UpdateProductRequest{},
			expectedStatus:        http.StatusBadRequest,
			expectedErrorResponse: "invalid id",
		},
		{
			name:                  "product not found",
			id:                    fmt.Sprintf("%d", 1111111111),
			body:                  model.UpdateProductRequest{},
			expectedStatus:        http.StatusNotFound,
			expectedErrorResponse: "product not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, err := json.Marshal(tt.body)
			require.NoError(t, err, "Failed to marshal request")

			url := fmt.Sprintf("/products/%s", tt.id)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", url, bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")

			ap := app.App()
			ap.ServeHTTP(w, req)

			var resp []byte
			if tt.expectedStatus == http.StatusOK {
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err, "Failed to unmarshal request")
			} else {
				resp, err = json.Marshal(tt.expectedErrorResponse)
				require.NoError(t, err, "Failed to unmarshal response")
			}

			require.Equal(t, tt.expectedStatus, w.Code)
			require.JSONEq(t, string(resp), w.Body.String())
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	//Prepare
	db := initDB(t)
	user, product, _ := createTestData(db, t, false)
	t.Cleanup(func() {
		removeData(db, user.ID)
	})

	testCases := []struct {
		name                  string
		id                    string
		expectedStatus        int
		expectedErrorResponse string
	}{
		{
			name:           "Success",
			id:             fmt.Sprintf("%d", product.ID),
			expectedStatus: http.StatusOK,
		},
		{
			name:                  "invalid id",
			id:                    "invalid",
			expectedStatus:        http.StatusBadRequest,
			expectedErrorResponse: "invalid id",
		},
		{
			name:                  "product not found",
			id:                    fmt.Sprintf("%d", 1111111111),
			expectedStatus:        http.StatusNotFound,
			expectedErrorResponse: "product not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			url := fmt.Sprintf("/products/%s", tt.id)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", url, nil)
			req.Header.Set("Content-Type", "application/json")

			ap := app.App()
			ap.ServeHTTP(w, req)

			var (
				resp []byte
				err  error
			)

			if tt.expectedStatus != http.StatusOK {
				resp, err = json.Marshal(tt.expectedErrorResponse)
				require.NoError(t, err, "Failed to unmarshal response")
			}

			require.Equal(t, tt.expectedStatus, w.Code)
			require.JSONEq(t, string(resp), w.Body.String()) // TODO check
		})
	}
}

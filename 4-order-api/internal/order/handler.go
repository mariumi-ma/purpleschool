package order

import (
	"errors"
	"net/http"
	"strconv"

	"purpleschool/internal/ctxutils"
	"purpleschool/internal/model"
	"purpleschool/internal/request"
	"purpleschool/internal/response"
)

type OrderHandler struct {
	OrderRepository *OrderRepository
	auth            AuthMiddleware
}

func NewOrderHandler(router *http.ServeMux, orderRepository *OrderRepository, auth AuthMiddleware) {
	handler := &OrderHandler{
		OrderRepository: orderRepository,
		auth:            auth,
	}

	router.Handle("GET /order/{id}", auth.IsAuth(handler.GetOrderByID()))
	router.Handle("GET /my-orders", auth.IsAuth(handler.GetOrderByUserID()))
	router.Handle("POST /order", auth.IsAuth(handler.CreateOrder()))
}

func (h *OrderHandler) GetOrderByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idPath := r.PathValue("id")

		id, err := strconv.Atoi(idPath)
		if err != nil {
			response.JSON(w, "invalid id", http.StatusBadRequest)
			return
		}

		userID := ctxutils.UserID(r.Context())

		order, err := h.OrderRepository.GetOrderByOrderIDAndUserID(uint(id), userID)
		if err != nil {
			if errors.Is(err, ErrOrderNotFound) {
				response.JSON(w, "order not found", http.StatusNotFound)
				return
			}
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, order, http.StatusOK)
	}
}

func (h *OrderHandler) GetOrderByUserID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := ctxutils.UserID(r.Context())

		orders, err := h.OrderRepository.GetOrdersByUserID(userID)
		if err != nil {
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, orders.ToResponse(), http.StatusOK)
	}
}

func (h *OrderHandler) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[model.CreateOrderRequest](w, r)
		if err != nil {
			response.JSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID := ctxutils.UserID(r.Context())

		order := &model.Order{}
		order.SetUserID(userID).
			SetProductsIDs(body.ProductIDs)

		order, err = h.OrderRepository.CreateOrder(order)
		if err != nil {
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := model.CreateOrderResponse{
			OrderID: order.ID,
		}

		response.JSON(w, resp, http.StatusCreated)
	}
}

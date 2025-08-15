package product

import (
	"errors"
	"net/http"
	"strconv"

	"purpleschool/internal/model"
	"purpleschool/internal/request"
	"purpleschool/internal/response"
)

type ProductHandler struct {
	ProductRepository *ProductRepository
}

func NewProductHandler(router *http.ServeMux, productRepository *ProductRepository) {
	handler := &ProductHandler{
		ProductRepository: productRepository,
	}

	router.Handle("GET /products/{id}", handler.GetProductByID())
	router.Handle("GET /products", handler.GetProducts())
	router.Handle("POST /products", handler.CreateProduct())
	router.Handle("PUT /products/{id}", handler.UpdateProduct())
	router.Handle("DELETE /products/{id}", handler.DeleteProduct())
}

func (h *ProductHandler) GetProductByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idPath := r.PathValue("id")

		id, err := strconv.Atoi(idPath)
		if err != nil {
			response.JSON(w, "invalid id", http.StatusBadRequest)
			return
		}

		product, err := h.ProductRepository.GetProductByID(uint(id))
		if err != nil {
			if errors.Is(err, ErrProductNotFound) {
				response.JSON(w, "product not found", http.StatusNotFound)
				return
			}
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, product, http.StatusOK)
	}
}

func (h *ProductHandler) GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := h.ProductRepository.GetProducts()
		if err != nil {
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, products.ToResponse(), http.StatusOK)
	}
}

func (h *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[model.CreateProductRequest](w, r)
		if err != nil {
			response.JSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		product := body.ToProduct()

		createdProduct, err := h.ProductRepository.CreateProduct(product)
		if err != nil {
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := createdProduct.ToResponse()
		response.JSON(w, resp, http.StatusCreated)
	}
}

func (h *ProductHandler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idPath := r.PathValue("id")

		id, err := strconv.Atoi(idPath)
		if err != nil {
			response.JSON(w, "invalid id", http.StatusBadRequest)
			return
		}

		body, err := request.HandleBody[model.UpdateProductRequest](w, r)
		if err != nil {
			return
		}

		product := body.ToProduct()
		product.SetID(int64(id))

		_, err = h.ProductRepository.GetProductByID(uint(id))
		if err != nil {
			if errors.Is(err, ErrProductNotFound) {
				response.JSON(w, "product not found", http.StatusNotFound)
				return
			}
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var updatedProduct *model.Product
		updatedProduct, err = h.ProductRepository.UpdateProduct(product)
		if err != nil {
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := updatedProduct.ToResponse()
		response.JSON(w, resp, http.StatusOK)
	}
}

func (h *ProductHandler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idPath := r.PathValue("id")

		id, err := strconv.Atoi(idPath)
		if err != nil {
			response.JSON(w, "invalid id", http.StatusBadRequest)
			return
		}

		_, err = h.ProductRepository.GetProductByID(uint(id))
		if err != nil {
			if errors.Is(err, ErrProductNotFound) {
				response.JSON(w, "product not found", http.StatusNotFound)
				return
			}
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.ProductRepository.DeleteProduct(uint(id))
		if err != nil {
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, nil, http.StatusOK)
	}
}

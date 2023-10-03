package handler

import (
	"fmt"
	"net/http"
	pb_product "preview_w2_p3/internal/product"

	"github.com/labstack/echo/v4"
)

// mungkin ada masalah pointer disini
type ClientHandler struct {
	ProductService pb_product.ProductServiceClient
}

func (ch ClientHandler) CreateProduct(c echo.Context) error {
	req := new(pb_product.AddProductRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request data")
	}

	product, err := ch.ProductService.AddProduct(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to add product: %v", err))
	}

	return c.JSON(http.StatusOK, product)
}
func (ch ClientHandler) ReadAllProduct(c echo.Context) error {
	resp, err := ch.ProductService.GetProduct(c.Request().Context(), &pb_product.GetProductRequest{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to get products: %v", err))
	}

	return c.JSON(http.StatusOK, resp.Products)
}
func (ch ClientHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")

	req := new(pb_product.UpdateProductRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request data")
	}
	req.Id = id

	product, err := ch.ProductService.UpdateProduct(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to update product: %v", err))
	}

	return c.JSON(http.StatusOK, product)
}
func (ch ClientHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")

	_, err := ch.ProductService.DeleteProduct(c.Request().Context(), &pb_product.DeleteProductRequest{Id: id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to delete product: %v", err))
	}

	return c.NoContent(http.StatusNoContent)
}

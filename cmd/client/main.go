package main

import (
	"log"
	"preview_w2_p3/handler"
	pb_product "preview_w2_p3/internal/product"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	ch := handler.ClientHandler{ProductService: pb_product.NewProductServiceClient(conn)}

	e := echo.New()

	p := e.Group("/product")
	p.POST("", ch.CreateProduct)
	p.GET("", ch.ReadAllProduct)
	p.PUT("/:id", ch.UpdateProduct)
	p.DELETE("/:id", ch.DeleteProduct)

	e.Logger.Fatal(e.Start(":8080"))
}

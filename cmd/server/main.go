package main

import (
	"context"
	"errors"
	"log"
	"net"
	pb "preview_w2_p3/internal/product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	products []*pb.Product
	pb.UnimplementedProductServiceServer
	Repo Repository
}

func (s *Server) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.Product, error) {
	newProduct := &pb.Product{}

	_, err := s.Repo.CreateProduct(req.Name, req.Description, req.Price, int(req.Stock))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}

	return newProduct, nil
}

func (s *Server) GetProduct(ctx context.Context, request *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	products, err := s.Repo.GetAllProduct()
	if err != nil {
		return nil, err
	}

	var productResponses []*pb.Product

	for _, product := range products {
		productResponse := &pb.Product{
			Id:          product.ID.Hex(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       int32(product.Stock),
		}
		productResponses = append(productResponses, productResponse)
	}

	response := &pb.GetProductResponse{
		Products: productResponses,
	}

	return response, nil
}

func (s *Server) UpdateProducts(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {

	result, err := s.Repo.UpdateProduct(req.Id, req.Name, req.Description, req.Price, int(req.Stock))
	if err != nil {
		return nil, err
	}

	updatedProduct := &pb.Product{
		Id:          result.ID.Hex(),
		Name:        result.Name,
		Description: result.Description,
		Price:       result.Price,
		Stock:       int32(result.Stock),
	}
	return updatedProduct, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	if req == nil || req.Id == "" {
		return nil, errors.New("Invalid request. Product ID is required.")
	}

	if err := s.Repo.DeleteProduct(req.Id); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func main() {

	db := ConnectDb()
	repo := NewRepository(db)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, &Server{Repo: repo})
	//protoc --go_out=. --go-grpc_out=. pb/user.proto

	log.Println("Server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC Server : %v", err)
	}
}

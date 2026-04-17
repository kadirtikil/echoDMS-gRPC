package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/echoDMS/db"
	"github.com/echoDMS/mtls"
	"github.com/echoDMS/proto/document"
	"github.com/echoDMS/proto/page"
	document_service "github.com/echoDMS/services/document"
	page_service "github.com/echoDMS/services/page"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	pool, err := db.NewPool(context.Background(), "postgresql://postgres:postgres@db:5432/echo?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer pool.Close()

	// gRPC SERVER
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	tlsCreds, err := mtls.LoadTLSCredentials()
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}
	fmt.Println(tlsCreds)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// REGISTER SERVICES HERE
	documentService := document_service.NewDocumentService(pool)
	pageService := page_service.NewPageService(pool)

	document.RegisterDocumentServiceServer(grpcServer, documentService)
	page.RegisterPageServiceServer(grpcServer, pageService)

	log.Println("gRPC server is running on: 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

}

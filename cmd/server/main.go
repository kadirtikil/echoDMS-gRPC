package main

import (
	"context"
	"log"
	"net"

	"github.com/echoCMS/db"
	"github.com/echoCMS/migrations"
	"google.golang.org/grpc"
)

func main() {
	pool, err := db.NewPool(context.Background(), "postgresql://postgres:postgres@db:5432/echo?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer pool.Close()

	if err := migrations.RunInitMigration(context.Background(), pool); err != nil {
		log.Fatalf("Failed to run initial migration: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	log.Println("gRPC server is running on: 50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

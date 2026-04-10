package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello, gRPC!")

	db, err := sql.Open("", "")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

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

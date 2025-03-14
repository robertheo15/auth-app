package main

import (
	"auth-app/internal/config"
	"auth-app/internal/repository"
	"auth-app/internal/service"
	"auth-app/pkg/proto/auth"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	config.LoadEnvFile()

	db, err := config.LoadDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	redis := config.NewRedis(ctx)
	repo := repository.NewUserRepository(db, redis)
	svc := service.NewUserService(repo)

	listener, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	auth.RegisterAuthServiceServer(server, svc)

	log.Println("gRPC server listening on :5001")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/Zagarazhi/go-project/generated"
	userservice "github.com/Zagarazhi/go-project/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Не забыть убрать
// protoc --go_out=generated --go_opt=paths=source_relative  --go-grpc_out=generated --go-grpc_opt=paths=source_relative  services.proto

func main() {
	addr := ":8080"

	apiServiceConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to ApiService: %v", err)
	}
	defer apiServiceConn.Close()

	apiServiceClient := generated.NewApiServiceClient(apiServiceConn)

	s := grpc.NewServer()
	userServiceServer := &userservice.UserServiceServer{
		ApiServiceClient: apiServiceClient,
	}

	generated.RegisterUserServiceServer(s, userServiceServer)

	fmt.Printf("HTTP Server is running on %s\n", addr)
	userservice.StartHTTPServer(userServiceServer, addr)
}

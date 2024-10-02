package main

import (
	// "context"
	// "fmt"
	// "log"
	// "net"
	"fmt"

	// "google.golang.org/grpc"

	c "github.com/hoang-cao-long/golang-side-projects/learn-gorm/config"
	// i "github.com/hoang-cao-long/golang-side-projects/learn-gorm/internal"
)

// // Define the gRPC server struct
// type Server struct{}

// // Implement the gRPC service methods
// func (s *Server) SayHello(ctx context.Context, request *HelloRequest) (*HelloResponse, error) {
// 	name := request.GetName()
// 	message := "Hello, " + name
// 	response := &HelloResponse{
// 		Message: message,
// 	}
// 	return response, nil
// }

func main() {
	// // Create a TCP listener on port 50051
	// listener, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	log.Fatalf("Failed to listen: %v", err)
	// }

	// // Create a new gRPC server
	// grpcServer := grpc.NewServer()

	// // Register the Server struct with the gRPC server
	// RegisterHelloServiceServer(grpcServer, &Server{})

	// // Start the gRPC server
	// if err := grpcServer.Serve(listener); err != nil {
	// 	log.Fatalf("Failed to serve: %v", err)
	// }

	// fmt.Print(i.Code)
	fmt.Print(c.Config{})

loop:
	for n := 0; n < 10; n++ {
		continue loop
	}
}

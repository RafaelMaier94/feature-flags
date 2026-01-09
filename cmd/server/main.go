package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/rafaelmaier/featureflags/proto/v1"
	"github.com/rafaelmaier/featureflags/repository"
	"github.com/rafaelmaier/featureflags/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" // ← ADD THIS LINE
)
func main() {
	port := 50051

	repo := repository.NewInMemoryRepository()
	log.Println("In-memory repo created")

	featureFlagService := service.NewFeatureFlagService(repo)
	log.Println("Service repo created")

	grpcServer := grpc.NewServer()
	log.Println("Server created")

	pb.RegisterFeatureAdminServiceServer(grpcServer, featureFlagService)
	log.Println("Feature admin service registered")
	
	reflection.Register(grpcServer)  // ← ADD THIS LINE
	log.Println("✓ gRPC reflection registered")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}
	log.Printf("✓ Server listening on port %d", port)

	go func() {
		if err := grpcServer.Serve(listener); err != nil{
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()

}
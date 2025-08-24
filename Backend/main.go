package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"fullstack-go-grpc/backend/controller"
	"fullstack-go-grpc/backend/service"
	"fullstack-go-grpc/database"
	pb "fullstack-go-grpc/protos/user"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = "50051"
	httpPort = "8081"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable not set")
	}

	// Create a context that is cancelled on interruption
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Database connection
	db := database.NewPostgresDB(dsn)
	defer db.Close()

	// Automatically create schema on startup
	if err := database.CreateSchema(ctx, db); err != nil {
		log.Fatalf("Failed to create database schema: %v", err)
	}

	// Initialize service and controller
	userService := service.NewUserService(db)
	userController := controller.NewUserController(userService)

	// Start gRPC server
	go func() {
		if err := runGrpcServer(userController); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start gRPC-Gateway (HTTP server)
	go func() {
		if err := runHttpServer(ctx); err != nil {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	log.Println("Backend service started...")
	<-ctx.Done() // Wait for interrupt signal
	log.Println("Shutting down servers...")
}

func runGrpcServer(userController *controller.UserController) error {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		return err
	}
	defer lis.Close()

	// Create a new gRPC server
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userController)
	reflection.Register(s) // Enable server reflection

	log.Printf("gRPC server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func runHttpServer(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	grpcEndpoint := "localhost:" + grpcPort
	log.Println("Connecting to gRPC server at", grpcEndpoint)

	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		return err
	}

	server := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	log.Printf("HTTP server listening at %s", server.Addr)
	return server.ListenAndServe()
}

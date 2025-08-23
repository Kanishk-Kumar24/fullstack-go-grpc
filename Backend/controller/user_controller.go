package controller

import (
	"context"
	"log"

	"fullstack-go-grpc/backend/service"
	"fullstack-go-grpc/database/internals/models"
	pb "fullstack-go-grpc/protos/user"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Three tier architecture?
// client --> server --> DB
// DB --> server --> client

// GRPC => g remot procdure call => function call/method call
// Communication / network layer ==> mux
// mux responsible hota hai, request ko sahi server ke method ke pass pauchana...

// how GRCP fast ???
// routes define .... 500 routes...
// network 500 routes...

// servers and server ke ander methods define kar diye
// 500 route grpc => 80-100 serives divide hoga...
// 5-6 methods honge harek service ke

// map[string][]string, servieName => ["method1", "method2"]
// O(1)
// ~1

// client/user ---> [grpcClient --> grpcServer] --> DB
// restAPI client --> server --> DB
// MVC => code structure...

// UserController implements the gRPC UserServiceServer interface.
type UserController struct {
	pb.UnimplementedUserServiceServer //
	userService                       service.UserService
}

// NewUserController creates a new UserController.
func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("CreateUser RPC called")
	var user models.User
	user.ConvertFromProto(req)

	createdUser, err := c.userService.CreateUser(ctx, &user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.CreateUserResponse{User: createdUser.ConvertToProto()}, nil
}

func (c *UserController) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("GetUser RPC called with ID: %s", req.UniqueId)
	id, err := uuid.Parse(req.UniqueId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID format: %v", err)
	}

	user, err := c.userService.GetUserByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.GetUserResponse{User: user.ConvertToProto()}, nil
}

func (c *UserController) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	log.Printf("UpdateUser RPC called with ID: %s", req.UniqueId)
	id, err := uuid.Parse(req.UniqueId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID format: %v", err)
	}

	// Fetch existing user
	user, err := c.userService.GetUserByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found for update: %v", err)
	}

	// Update fields
	user.Name = req.Name
	user.PhoneNumber = req.PhoneNumber
	if req.Address != nil {
		user.Country = req.Address.Country
		user.State = req.Address.State
	}

	updatedUser, err := c.userService.UpdateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.UpdateUserResponse{User: updatedUser.ConvertToProto()}, nil
}

func (c *UserController) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	log.Printf("DeleteUser RPC called with ID: %s", req.UniqueId)
	id, err := uuid.Parse(req.UniqueId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID format: %v", err)
	}

	err = c.userService.DeleteUser(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUserResponse{Message: "User deleted successfully"}, nil
}

func (c *UserController) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	log.Printf("ListUsers RPC called")
	users, err := c.userService.ListUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, user.ConvertToProto())
	}

	return &pb.ListUsersResponse{Users: pbUsers}, nil
}

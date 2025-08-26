package controller

import (
	"context"
	"log"
	"fullstack-go-grpc/backend/service"
	pb "fullstack-go-grpc/protos/user"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserController struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log.Printf("CreateUser RPC called")
	createdUser, err := c.userService.CreateUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	return createdUser, nil
}

func (c *UserController) GetUser(ctx context.Context, req *pb.U_ID) (*pb.UserResponse, error) {
	log.Printf("GetUser RPC called with ID: %s", req.UniqueId)
	id, err := uuid.Parse(req.UniqueId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID format: %v", err)
	}
	user, err := c.userService.GetUserByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}
	return &pb.UserResponse{User: user.ConvertToProto()}, nil
}

func (c *UserController) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	log.Printf("UpdateUser RPC called with ID: %s", req.UniqueId)
	updatedUser, err := c.userService.UpdateUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}
	return updatedUser, nil
}

func (c *UserController) DeleteUser(ctx context.Context, req *pb.U_ID) (*pb.DeleteUserResponse, error) {
	log.Printf("DeleteUser RPC called with ID: %s", req.UniqueId)
	deletedUser, err := c.userService.DeleteUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return deletedUser, nil
}

func (c *UserController) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	log.Printf("ListUsers RPC called")
	users, err := c.userService.ListUsers(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	return users, nil
}
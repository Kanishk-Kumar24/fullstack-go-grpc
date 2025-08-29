package service

import (
	"context"
	"fullstack-go-grpc/internals/models"
	pb "fullstack-go-grpc/protos/user"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"fullstack-go-grpc/backend/repo"
	"time"
)

// UserService provides business logic for user operations.
// type UserService interface {
// 	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*models.User, error)
// 	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
// 	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*models.User, error)
// 	DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) error
// 	ListUsers(ctx context.Context) ([]models.User, error)
// }


type UserService struct{
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) *UserService{
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	var user models.User
	user.ConvertFromProto(req.User)
	err := s.userRepo.CreateUser(ctx,&user)
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{User: user.ConvertToProto()}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.userRepo.GetUserByID(ctx,id,&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	
	id,err:= uuid.Parse(req.UniqueId)
	if err!=nil{
		return nil,err
	}
	user,err:= s.GetUserByID(ctx,id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found for update: %v", err)
	}
	user.Name = req.Name
	user.PhoneNumber = req.PhoneNumber
	if req.Address != nil {
		user.Country = req.Address.Country
		user.State = req.Address.State
	}
	user.UpdatedAt = time.Now()
	er := s.userRepo.UpdateUser(ctx,user)
	if er != nil {
		return nil, er
	}
	return &pb.UserResponse{User: user.ConvertToProto()}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.UserGetterRequest) (*pb.DeleteUserResponse, error) {
	id,err:=uuid.Parse(req.UniqueId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID format: %v", err)
	}
	err = s.userRepo.DeleteUser(ctx,id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{Message: "User deleted successfully"}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var users []models.User
	err := s.userRepo.ListUsers(ctx, &users)
	if err != nil {
		return nil, err
	}
	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, user.ConvertToProto())
	}

	return &pb.ListUsersResponse{Users: pbUsers}, nil
}
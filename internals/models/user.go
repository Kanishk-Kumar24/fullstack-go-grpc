package models

import (
	"context"
	"time"

	pb "fullstack-go-grpc/protos/user"

	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// User represents the user model in the database.
type User struct {
	BaseModel `bun:"table:users,alias:u"`

	Name        string    `bun:"name,notnull"`
	Email       string    `bun:"email,unique,notnull"`
	PhoneNumber string    `bun:"phone_number"`
	DOB         time.Time `bun:"dob,nullzero"`
	Country     string    `bun:"country"`
	State       string    `bun:"state"`
}

// BeforeAppendModel is a Bun hook that sets timestamps before creating.
func (u *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
	return nil
}

// BeforeUpdateModel is a Bun hook that sets the 'updated_at' timestamp before updating.
func (u *User) BeforeUpdateModel(ctx context.Context, query bun.Query) error {
	u.UpdatedAt = time.Now()
	return nil
}

// ConvertToProto converts a database User model to a gRPC User message.
func (u *User) ConvertToProto() *pb.User {
	user := &pb.User{
		UniqueId:    u.UniqueID.String(),
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Address: &pb.Address{
			Country: u.Country,
			State:   u.State,
		},
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
	if !u.DOB.IsZero() {
		user.Dob = timestamppb.New(u.DOB)
	}
	return user
}

// ConvertFromProto converts a gRPC CreateUserRequest to a database User model.
func (u *User) ConvertFromProto(req *pb.User) {
	u.Name = req.Name
	u.Email = req.Email
	u.PhoneNumber = req.PhoneNumber
	if req.Dob != nil {
		u.DOB = req.Dob.AsTime()
	}
	if req.Address != nil {
		u.Country = req.Address.Country
		u.State = req.Address.State
	}
}

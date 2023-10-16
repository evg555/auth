package converter

import (
	"github.com/evg555/auth/internal/model"
	proto "github.com/evg555/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToGetResponseFromService(user *model.User) *proto.GetResponse {
	var updatedAt *timestamppb.Timestamp

	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &proto.GetResponse{
		Id:        user.ID,
		Name:      user.Name.String,
		Email:     user.Email.String,
		Role:      proto.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserFromCreateRequest(req *proto.CreateRequest) *model.User {
	user := &model.User{
		Password: req.GetPassword(),
		Role:     int32(req.Role),
	}

	user.Name.String = req.GetName()
	user.Email.String = req.GetEmail()

	return user
}

func ToUserFromUpdateRequest(req *proto.UpdateRequest) *model.User {
	var user model.User

	if req.GetName() != nil {
		user.Name.String = req.GetName().Value
		user.Name.Valid = true
	}

	if req.GetEmail() != nil {
		user.Email.String = req.GetEmail().Value
		user.Email.Valid = true
	}

	user.ID = req.GetId()

	return &user
}

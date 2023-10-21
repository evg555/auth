package converter

import (
	"github.com/evg555/auth/internal/model"
	modelRepo "github.com/evg555/auth/internal/repository/user/model"
)

func ToUserFromRepo(user modelRepo.User) *model.User {
	modelUser := model.User{
		ID:        user.ID,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	modelUser.Name.String = user.Name
	modelUser.Email.String = user.Name

	return &modelUser
}

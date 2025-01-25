package user

import (
	"backend/models"

	"github.com/google/uuid"
)

type userTokenResponse struct {
	User struct {
		Id       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		Token    string    `json:"token"`
	} `json:"user"`
}

func newUserTokenResponse(u *models.User, token string) *userTokenResponse {
	r := &userTokenResponse{}
	r.User.Username = u.Username
	r.User.Id = u.Id
	r.User.Token = token
	return r
}

type userResponse struct {
	User struct {
		Id       uuid.UUID `json:"id"`
		Username string    `json:"username"`
	} `json:"user"`
}

func newUserResponse(u *models.User) *userResponse {
	r := &userResponse{}
	r.User.Username = u.Username
	r.User.Id = u.Id
	return r
}

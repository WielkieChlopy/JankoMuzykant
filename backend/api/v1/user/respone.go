package user

import (
	"backend/models"

	"github.com/google/uuid"
)

type userTokenResponse struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Token    string    `json:"token"`
}

func newUserTokenResponse(u *models.User, token string) *userTokenResponse {
	r := &userTokenResponse{}
	r.Username = u.Username
	r.Id = u.Id
	r.Token = token
	return r
}

type userResponse struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func newUserResponse(u *models.User) *userResponse {
	r := &userResponse{}
	r.Username = u.Username
	r.Id = u.Id
	return r
}

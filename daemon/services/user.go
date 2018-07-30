package services

import (
	"github.com/yankeguo/bunker/types"
)

// UserService user service of rpc
type UserService struct {
	Service
}

// GetUsers get all users from database
func (s *UserService) GetUsers() (err error) {
	return
}

// CreateUser create a user
func (s *UserService) CreateUser(req types.CreateUserRequest, out *int) (err error) {
	return
}

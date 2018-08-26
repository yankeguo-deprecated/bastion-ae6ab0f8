package models

import (
	"github.com/jinzhu/copier"
	"github.com/yankeguo/bastion/types"
)

// User user model in boltdb, mirror of types.User
type User struct {
	Account        string `storm:"id"` // primary key
	Nickname       string
	PasswordDigest string
	PasswordFailed int64
	IsAdmin        bool
	IsBlocked      bool
	CreatedAt      int64
	UpdatedAt      int64
	ViewedAt       int64
}

func (u User) ToGRPCUser() *types.User {
	n := types.User{}
	copier.Copy(&n, &u)
	return &n
}

package db

import (
	"time"

	"github.com/yankeguo/bastion/types"
)

// User user model
type User struct {
	Model
	Account        string     `gorm:"not null;unique_index"`
	Nickname       string     `gorm:"not null;type:text"`
	PasswordDigest string     `gorm:"not null;type:text"`
	IsAdmin        int        `gorm:"not null;default:0"`
	IsBlocked      int        `gorm:"not null;default:0"`
	UsedAt         *time.Time `gorm:""`
}

// ToRPCUser to rpc user
func (u User) ToRPCUser() (g *types.User) {
	g = &types.User{}
	g.Account = u.Account
	g.Nickname = u.Nickname
	g.IsAdmin = u.IsAdmin != 0
	g.IsBlocked = u.IsBlocked != 0
	g.CreatedAt = u.CreatedAt.Unix()
	g.UpdatedAt = u.UpdatedAt.Unix()
	if u.UsedAt != nil {
		g.ViewedAt = u.UsedAt.Unix()
	}
	return
}

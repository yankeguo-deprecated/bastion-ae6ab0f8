package db

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/yankeguo/bastion/types"
)

// User user model
type User struct {
	Model
	Account        string     `gorm:"not null;unique_index" json:"account"`    // account name
	Nickname       string     `gorm:"not null;type:text" json:"nickname"`      // nickname
	PasswordDigest string     `gorm:"not null;type:text" json:"-"`             // password encrypted by bcrypt
	IsAdmin        bool       `gorm:"not null;default:false" json:"isAdmin"`   // is this user system admin
	IsBlocked      bool       `gorm:"not null;default:false" json:"isBlocked"` // is this user blocked
	UsedAt         *time.Time `gorm:"" json:"usedAt"`                          // last seen at
}

// BeforeSave before save callback
func (u *User) BeforeSave() (err error) {
	if !NamePattern.MatchString(u.Account) {
		err = errors.New(`invalid field user.account, allows 3~15 letters, numbers, "_" or "-"`)
	}
	return
}

// SetPassword update password for user
// bcrypt produces clear text encrypted password, no further encoding needed
func (u *User) SetPassword(p string) (err error) {
	var b []byte
	if b, err = bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost); err != nil {
		return
	}
	u.PasswordDigest = string(b)
	return
}

// CheckPassword check password
func (u User) CheckPassword(p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(p)) == nil
}

// ToRPCUser to rpc user
func (u User) ToRPCUser() (g *types.User) {
	g = &types.User{}
	g.Account = u.Account
	g.Nickname = u.Nickname
	g.IsAdmin = u.IsAdmin
	g.IsBlocked = u.IsBlocked
	g.CreatedAt = u.CreatedAt.Unix()
	g.UpdatedAt = u.UpdatedAt.Unix()
	if u.UsedAt != nil {
		g.ViewedAt = u.UsedAt.Unix()
	}
	return
}

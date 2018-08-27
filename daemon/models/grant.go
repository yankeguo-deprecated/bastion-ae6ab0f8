package models

import (
	"github.com/jinzhu/copier"
	"github.com/yankeguo/bastion/types"
)

type Grant struct {
	Id              string `storm:"id"`
	Account         string `storm:"index"`
	HostnamePattern string
	User            string
	ExpiredAt       int64
	CreatedAt       int64
}

func (n Grant) BuildId() string {
	return n.Account + "$" + n.HostnamePattern + "$" + n.User
}

func (n Grant) ToGRPCGrant() *types.Grant {
	v := types.Grant{}
	copier.Copy(&v, &n)
	return &v
}

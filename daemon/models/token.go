package models

import (
	"github.com/jinzhu/copier"
	"github.com/yankeguo/bastion/types"
)

type Token struct {
	Id          int64  `storm:"id,increment"`
	Token       string `storm:"unique"`
	Account     string `storm:"index"`
	Description string
	CreatedAt   int64
	ViewedAt    int64
}

func (n Token) ToGRPCToken() *types.Token {
	o := types.Token{}
	copier.Copy(&o, &n)
	return &o
}

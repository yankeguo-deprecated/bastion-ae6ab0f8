package models

import (
	"github.com/jinzhu/copier"
	"github.com/yankeguo/bastion/types"
)

// Node server model
type Node struct {
	Hostname  string `storm:"id"`
	User      string
	Address   string
	Source    string
	CreatedAt int64
}

func (n Node) ToGRPCNode() *types.Node {
	o := types.Node{}
	copier.Copy(&o, &n)
	return &o
}

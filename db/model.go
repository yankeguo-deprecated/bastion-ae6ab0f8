package db

import "time"

// Model basic model, not using orm.Model, no deletedAt
type Model struct {
	ID        uint      `orm:"primary_key" json:"id"` // id
	CreatedAt time.Time `orm:"" json:"createdAt"`     // created at
	UpdatedAt time.Time `orm:"" json:"updatedAt"`     // updated at
}

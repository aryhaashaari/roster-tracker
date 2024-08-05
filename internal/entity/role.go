// Package entity
// Automatic generated
package entity

import (
	"time"
)

// Role entity
type Role struct{
	Code string	`db:"code,omitempty" json:"code"`
	Name string	`db:"name,omitempty" json:"name"`
	CreatedAt time.Time	`db:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time	`db:"updated_at,omitempty" json:"updated_at"`
	DeletedAt *time.Time	`db:"deleted_at,omitempty" json:"deleted_at"`
}

package model

import "time"

type User struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	Roles      []Role    `json:"roles"`
	LastAccess time.Time `json:"last_access"`
}

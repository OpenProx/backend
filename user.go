package backend

import "time"

// CreateUser creates a new user
func CreateUser() *User {
	return &User{
		Joined:     time.Now().Unix(),
		LastActive: time.Now().Unix(),
	}
}

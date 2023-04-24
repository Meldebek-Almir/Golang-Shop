package models

import "time"

type User struct {
	ID        int       `json:"user_id"`
	Nickname  string    `json:"nickname"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	Created   time.Time `json:"created_time"`
	Updated   time.Time `json:"updated_time"`
}

type UserDataProduct struct {
	User    User      `json:"user"`
	Product []Product `json:"product"`
}

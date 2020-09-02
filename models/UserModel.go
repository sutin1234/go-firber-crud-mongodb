package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User struct DB:users
type User struct {
	gorm.Model
	ID       string
	Name     string
	LName    string
	Age      int
	Email    string
	UserName string
	Password string
	Token    string
	Birthday time.Time
	CreateAt time.Time
	UpdateAt time.Time
}

// Users type
type Users []User

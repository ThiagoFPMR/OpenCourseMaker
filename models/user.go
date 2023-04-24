package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nome     string
	Email    string
	Password string
	pswdhash string
	Active   string
	timeout  string
}

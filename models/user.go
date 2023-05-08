package models

type User struct {
	ID       uint   `gorm:"primary_Key"`
	Nome     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Tipo     int    `gorm:"not null"`
	Pswdhash string `gorm:"not null"`
}

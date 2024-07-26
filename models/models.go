package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"not Null"`
	Email    string `gorm:"not Null;unique"`
	Password string `gorm:"not Null"`
}
type Admin struct {
	Name     string
	Email    string
	Password string
}

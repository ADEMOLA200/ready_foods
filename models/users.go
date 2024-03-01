package models

import "time"

type User struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	FirstName       string `json:"firstname" gorm:"not null"`
	LastName        string `json:"lastname" gorm:"not null"`
	Username        string `json:"username" gorm:"unique;type:varchar(100);not null"`
	Email           string `json:"email" gorm:"unique;type:varchar(191)"`
	Password        string `json:"password" gorm:"not null"`
	ConfirmPassword string `json:"confirm-password" gorm:"-"`
}

type RevokedToken struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Token     string `gorm:"type:varchar(255);unique_index" json:"token"`
	RevokedAt time.Time
}
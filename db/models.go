package db

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username    string `gorm:"varchar(50),unique"`
	Email       string `gorm:"varchar(25), unique"`
	Password    string `gorm:"varchar(64)"`
	FirstName   string `gorm:"varchar(50)"`
	LastName    string `gorm:"varchar(50)"`
	PhoneNumber string `gorm:"varchar(30),unique"`
	Gender      string `gorm:"type:varchar(50)"`
}

type Book struct {
	gorm.Model
	Name              string    `gorm:"type:varchar(255)"`
	Category          string    `gorm:"type:varchar(255)"`
	Volume            int       `gorm:"type:integer"`
	PublishedAt       time.Time `gorm:"type:date"`
	Summary           string    `gorm:"type:text"`
	Publisher         string    `gorm:"type:varchar(255)"`
	AuthorFirstName   string    `gorm:"type:varchar(50)"`
	AuthorLastName    string    `gorm:"type:varchar(50)"`
	AuthorBirthday    time.Time `gorm:"type:date"`
	AuthorNationality string    `gorm:"type:varchar(50)"`
	UserID            uint
}

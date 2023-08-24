package db

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username    string `gorm:"varchar(50),unique" json:"username"`
	Email       string `gorm:"varchar(25), unique" json:"email"`
	Password    string `gorm:"varchar(64)" json:"password"`
	FirstName   string `gorm:"varchar(50)" json:"first_name"`
	LastName    string `gorm:"varchar(50)" json:"last_name"`
	PhoneNumber string `gorm:"varchar(30),unique" json:"phone_number"`
	Gender      string `gorm:"type:varchar(50)" json:"gender"`
}

type Book struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(255)"`
	Category    string    `gorm:"type:varchar(255)"`
	Volume      int       `gorm:"type:integer"`
	PublishedAt time.Time `gorm:"type:date"`
	Summary     string    `gorm:"type:text"`
	Publisher   string    `gorm:"type:varchar(255)"`
	FirstName   string    `gorm:"type:varchar(50)"`
	LastName    string    `gorm:"type:varchar(50)"`
	Birthday    time.Time `gorm:"type:date"`
	Nationality string    `gorm:"type:varchar(50)"`
	UserID      uint
}

package db

import (
	"gorm.io/gorm"
	"time"
)

// User struct which holds user info
type User struct {
	gorm.Model
	Username    string `gorm:"varchar(50),unique" json:"user_name"`
	Email       string `gorm:"varchar(25), unique" json:"email"`
	Password    string `gorm:"varchar(64)" json:"password"`
	FirstName   string `gorm:"varchar(50)" json:"first_name"`
	LastName    string `gorm:"varchar(50)" json:"last_name"`
	PhoneNumber string `gorm:"varchar(30),unique" json:"phone_number"`
	Gender      string `gorm:"type:varchar(50)" json:"gender"`
}

// Book struct which holds book info
type Book struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	Category    string    `gorm:"type:varchar(255)" json:"category"`
	Volume      int       `gorm:"type:integer" json:"volume"`
	PublishedAt time.Time `gorm:"type:date" json:"published_at"`
	Summary     string    `gorm:"type:text" json:"summary"`
	Publisher   string    `gorm:"type:varchar(255)" json:"publisher"`
	Author      Author    `gorm:"embedded" json:"author"`
	UserID      uint      `gorm:"type:int"`
}

type Author struct {
	FirstName   string    `gorm:"type:varchar(50)" json:"first_name"`
	LastName    string    `gorm:"type:varchar(50)" json:"last_name"`
	Birthday    time.Time `gorm:"type:date" json:"birthday"`
	Nationality string    `gorm:"type:varchar(50)" json:"nationality"`
}

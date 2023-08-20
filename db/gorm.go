package db

import (
	"bookman/config"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DuplicateEmailError       = errors.New("This Email is Already Taken")
	DuplicateUsernameError    = errors.New("This Username is Already Taken")
	DuplicatePhoneNumberError = errors.New("This Phone Number is Already Taken")
)

type GormDB struct {
	cfg config.Config
	db  *gorm.DB
}

func NewGormDB(cfg config.Config) (*GormDB, error) {
	c := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.Username,
		cfg.Database.Password,
	)

	db, err := gorm.Open(postgres.Open(c), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &GormDB{
		cfg,
		db,
	}, nil

}

func (gdb *GormDB) CreateSchemas() (error, error) {
	err1 := gdb.db.AutoMigrate(&User{})
	err2 := gdb.db.AutoMigrate(&Book{})
	if err1 != nil || err2 != nil {
		return err1, err2
	}
	return nil, nil
}

func (gdb *GormDB) CreateNewUser(user *User) error {
	var count int64
	gdb.db.Model(&User{}).Where("username = ?", user.Username).Count(&count)
	if count != 0 {
		return DuplicateUsernameError
	}

	gdb.db.Model(&User{}).Where("email = ?", user.Email).Count(&count)
	if count != 0 {
		return DuplicateEmailError
	}

	gdb.db.Model(&User{}).Where("phone_number = ?", user.PhoneNumber).Count(&count)
	if count != 0 {
		return DuplicatePhoneNumberError
	}

	EncryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if err != nil {
		return err
	}

	user.Password = string(EncryptedPassword)

	response := gdb.db.Create(user)
	return response.Error
}

func (gdb *GormDB) GetUserByUsername(username string) (*User, error) {
	var user User
	err := gdb.db.Where(&User{Username: username}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

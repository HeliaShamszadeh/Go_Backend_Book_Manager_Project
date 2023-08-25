package db

import (
	"bookman/config"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// defining possible errors for defining database and CRUD
var (
	DuplicateEmailError       = errors.New("This Email is Already Taken")
	DuplicateUsernameError    = errors.New("This Username is Already Taken")
	DuplicatePhoneNumberError = errors.New("This Phone Number is Already Taken")
	GenderNotAllowedError     = errors.New("Only female, male, or others are acceptable as genders")
	UserNameNotFoundError     = errors.New("User not found")
)

// GormDB is a struct which keeps info of config and database
type GormDB struct {
	cfg config.Config
	Db  *gorm.DB
}

// NewGormDB creates a new database (as a struct)
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

// CreateSchemas creates the tables needed for keeping records(Users and Books)
func (gdb *GormDB) CreateSchemas() (error, error) {
	err1 := gdb.Db.AutoMigrate(&User{})
	err2 := gdb.Db.AutoMigrate(&Book{})
	if err1 != nil || err2 != nil {
		return err1, err2
	}
	return nil, nil
}

// CreateNewUser inserts a new user to Users' table
func (gdb *GormDB) CreateNewUser(user *User) error {
	var count int64
	gdb.Db.Model(&User{}).Where("username = ?", user.Username).Count(&count)
	if count != 0 {
		return DuplicateUsernameError
	}

	gdb.Db.Model(&User{}).Where("email = ?", user.Email).Count(&count)
	if count != 0 {
		return DuplicateEmailError
	}

	gdb.Db.Model(&User{}).Where("phone_number = ?", user.PhoneNumber).Count(&count)
	if count != 0 {
		return DuplicatePhoneNumberError
	}

	if !(user.Gender == "male" || user.Gender == "female" || user.Gender == "others") {
		return GenderNotAllowedError
	}
	EncryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if err != nil {
		return err
	}

	user.Password = string(EncryptedPassword)

	response := gdb.Db.Create(user)
	return response.Error
}

// GetUserByUsername finds a record by selecting on username
func (gdb *GormDB) GetUserByUsername(username string) (*User, error) {
	var user User
	err := gdb.Db.Where(&User{Username: username}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateNewBook inserts new record of book to Books' table
func (gdb *GormDB) CreateNewBook(book *Book) error {
	return gdb.Db.Create(&book).Error
}

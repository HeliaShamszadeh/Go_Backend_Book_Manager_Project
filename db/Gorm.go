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
	DuplicateEmailErr       = errors.New("this email is already taken")
	DuplicateUsernameErr    = errors.New("this username is already taken")
	DuplicatePhoneNumberErr = errors.New("this phone number is already taken")
	GenderNotAllowedErr     = errors.New("only female, male, or others are acceptable as genders")
	UserNameNotFoundErr     = errors.New("user not found")
	BookNotFoundErr         = errors.New("book not found")
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
func (gdb *GormDB) CreateSchemas() (error, error, error) {
	err1 := gdb.Db.AutoMigrate(&User{})
	err2 := gdb.Db.AutoMigrate(&Book{})
	err3 := gdb.Db.AutoMigrate(&Content{})
	if err1 != nil || err2 != nil || err3 != nil {
		return err1, err2, err3
	}
	return nil, nil, nil
}

// CreateNewUser inserts a new user to Users' table
func (gdb *GormDB) CreateNewUser(user *User) error {
	var count int64
	gdb.Db.Model(&User{}).Where("username = ?", user.Username).Count(&count)
	if count != 0 {
		return DuplicateUsernameErr
	}

	gdb.Db.Model(&User{}).Where("email = ?", user.Email).Count(&count)
	if count != 0 {
		return DuplicateEmailErr
	}

	gdb.Db.Model(&User{}).Where("phone_number = ?", user.PhoneNumber).Count(&count)
	if count != 0 {
		return DuplicatePhoneNumberErr
	}

	if !(user.Gender == "male" || user.Gender == "female" || user.Gender == "others") {
		return GenderNotAllowedErr
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

// GetAllBooks retrieves all records saved in books table
func (gdb *GormDB) GetAllBooks() (*[]Book, error) {
	var books []Book
	err := gdb.Db.Model(Book{}).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return &books, nil

}

// GetBookById finds a book in books table by selecting on id column
func (gdb *GormDB) GetBookById(id int) (*Book, error) {
	var count int64
	gdb.Db.Model(&Book{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return nil, BookNotFoundErr
	}
	var book Book
	err := gdb.Db.Where("id = ?", id).First(&book).Error
	if err != nil {
		return nil, err
	} else {
		return &book, nil
	}
}

// GetBookContents gets contents of a book
func (gdb *GormDB) GetBookContents(id int) (*[]Content, error) {
	var contents []Content
	err := gdb.Db.Model(Content{}).Where("book_id = ?", id).Find(&contents).Error
	if err != nil {
		return nil, err
	}
	return &contents, nil
}

// UpdateBook updates name and category of a book
func (gdb *GormDB) UpdateBook(book *Book, name, category string) error {
	var count int64
	gdb.Db.Model(&Book{}).Where("id = ?", book.ID).Count(&count)
	if count == 0 {
		return BookNotFoundErr
	}
	return gdb.Db.Model(Book{}).Where("id = ?", book.ID).Update("name", name).Update("category", category).Error
}

// DeleteBookById deletes a book based on the given id
func (gdb *GormDB) DeleteBookById(id int) error {
	var count int64
	gdb.Db.Model(&Book{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return BookNotFoundErr
	}
	return gdb.Db.Delete(&Book{}, id).Error
}

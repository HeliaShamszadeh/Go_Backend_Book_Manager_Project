package handler

import (
	"bookman/db"
	"time"
)

// LoginRequestBody struct for LoginHandler
type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateBookRequestBody for CreateBookHandler
type CreateBookRequestBody struct {
	Name            string    `json:"name"`
	Author          author    `json:"author"`
	Category        string    `json:"category"`
	Volume          int       `json:"volume"`
	PublishedAt     time.Time `json:"published_at"`
	Summary         string    `json:"summary"`
	TableOfContents []string  `json:"table_of_contents"`
	Publisher       string    `json:"publisher"`
}

type author struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Birthday    time.Time `json:"birthday"`
	Nationality string    `json:"nationality"`
}

// Books struct for GetAllBooksHandler
type Books struct {
	Books *[]db.Book
}

// GetBooksResponseBody for GetAllBooksHandler and GetBookByIdHandler
type GetBooksResponseBody struct {
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Category    string    `json:"category"`
	Volume      int       `json:"volume"`
	PublishedAt time.Time `json:"published_at"`
	Summary     string    `json:"summary"`
	Publisher   string    `json:"publisher"`
}

// UpdateRequestBody Struct for UpdateBookHandler
type UpdateRequestBody struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

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

// Books struct for GetAllBooksHandler
type Books struct {
	Books *[]db.Book
}

// GetAllBooksResponseBody for GetAllBooksHandler and GetBookByIdHandler
type GetAllBooksResponseBody struct {
	BookId      uint      `json:"id"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Category    string    `json:"category"`
	Volume      int       `json:"volume"`
	PublishedAt time.Time `json:"published_at"`
	Summary     string    `json:"summary"`
	Publisher   string    `json:"publisher"`
}

// GetBookByIdResponseBody for GetBookByIdHandler
type GetBookByIdResponseBody struct {
	Name            string    `json:"name"`
	Author          string    `json:"author"`
	Category        string    `json:"category"`
	Volume          int       `json:"volume"`
	PublishedAt     time.Time `json:"published_at"`
	Summary         string    `json:"summary"`
	TableOfContents []string  `json:"table_of_contents"`
	Publisher       string    `json:"publisher"`
}

// UpdateRequestBody Struct for UpdateBookHandler
type UpdateRequestBody struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

// TableOfContents for GetBookByIdHandler
type TableOfContents struct {
	Contents []string `json:"table_of_contents"`
}

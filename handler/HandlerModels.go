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

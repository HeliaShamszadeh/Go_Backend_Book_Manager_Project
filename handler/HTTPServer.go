package handler

import (
	"bookman/authenticate"
	"bookman/db"
	"github.com/sirupsen/logrus"
	"net/http"
)

type BookManagerServer struct {
	DB           *db.GormDB
	Logger       *logrus.Logger
	Authenticate *authenticate.Authenticate
}

// BooksRootHandler sets the final handler function based on request method for /api/v1/books
func (bm *BookManagerServer) BooksRootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		bm.CreateBookHandler(w, r)
	} else if r.Method == http.MethodGet {
		bm.GetAllBooksHandler(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// BooksSubTreeHandler sets the final handler function based on request method for /api/v1/books/<id>
func (bm *BookManagerServer) BooksSubTreeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		bm.GetBookByIdHandler(w, r)
	} else if r.Method == http.MethodPut {
		bm.UpdateBookHandler(w, r)
	} else if r.Method == http.MethodDelete {
		bm.DeleteBookHandler(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

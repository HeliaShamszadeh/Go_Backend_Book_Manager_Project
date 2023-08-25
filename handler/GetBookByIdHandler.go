package handler

import (
	"bookman/authenticate"
	"bookman/db"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
)

func (bm *BookManagerServer) GetBookByIdHandler(w http.ResponseWriter, r *http.Request) {
	// check the request method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get access token from header and check related errors
	AuthorizationToken := r.Header.Get("Authorization")
	if AuthorizationToken == "" {
		response := map[string]interface{}{
			"message": authenticate.InvalidTokenErr,
		}
		resBody, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resBody)
		return
	}

	// check if this username exists
	_, err := bm.Authenticate.GetUsernameByToken(AuthorizationToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bm.Logger.WithError(err).Warn("could not login the user")
		return
	}

	// get book id from URL
	BookIdStr := path.Base(r.URL.Path)
	if BookIdStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// parse BookId to int
	BookIdInt, err := strconv.Atoi(BookIdStr)
	if err != nil {
		bm.Logger.WithError(err).Warn("error occurred while parsing book id")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// retrieve the book from database
	ReturnedBook, err := bm.DB.GetBookById(BookIdInt)
	if err != nil {
		if err == db.BookNotFoundError {
			w.WriteHeader(http.StatusBadRequest)
			resBody, _ := json.Marshal(map[string]interface{}{"message": "book not found"})
			w.Write(resBody)
			return
		}
		bm.Logger.WithError(err).Warn("error retrieving book from database (GetBookById)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// append exported data to a preferred data structure for marshalling(JSON)
	Book := GetBooksResponseBody{
		Name:        ReturnedBook.Name,
		Author:      ReturnedBook.FirstName + " " + ReturnedBook.LastName,
		Category:    ReturnedBook.Category,
		Volume:      ReturnedBook.Volume,
		PublishedAt: ReturnedBook.PublishedAt,
		Summary:     ReturnedBook.Summary,
		Publisher:   ReturnedBook.Publisher,
	}

	// marshal the data
	resBody, err := json.Marshal(Book)
	if err != nil {
		bm.Logger.WithError(err).Warn("error marshalling data (GetBookById)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBody)

}

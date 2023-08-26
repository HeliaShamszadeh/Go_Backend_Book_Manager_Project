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
		resBody, _ := json.Marshal(map[string]interface{}{
			"message": authenticate.EmptyTokenErr.Error(),
		})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resBody)
		return
	}

	// check if this username exists
	_, err := bm.Authenticate.GetUsernameByToken(AuthorizationToken)
	if err != nil {
		if err == authenticate.CannotValidateTokenErr {
			bm.Logger.WithError(err).Warn()
			w.WriteHeader(http.StatusBadRequest)
			response, _ := json.Marshal(map[string]interface{}{"message": authenticate.CannotValidateTokenErr.Error()})
			w.Write(response)
			return
		} else if err == authenticate.InvalidTokenErr {
			bm.Logger.WithError(err).Warn()
			w.WriteHeader(http.StatusBadRequest)
			response, _ := json.Marshal(map[string]interface{}{"message": authenticate.InvalidTokenErr.Error()})
			w.Write(response)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			bm.Logger.WithError(err).Warn("there was a problem with access token - maybe user does not exist")
			return
		}
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
		if err == db.BookNotFoundErr {
			w.WriteHeader(http.StatusBadRequest)
			resBody, _ := json.Marshal(map[string]interface{}{"message": "book not found"})
			w.Write(resBody)
			return
		}
		bm.Logger.WithError(err).Warn("error retrieving book from database (GetBookById)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get book's contents
	Contents, err := bm.DB.GetBookContents(BookIdInt)
	if err != nil {
		bm.Logger.WithError(err).Warn("error retrieving book's contents from database (GetBookContents)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// append exported data to a preferred data structure for marshalling(JSON)
	Book := GetBookByIdResponseBody{
		Name:        ReturnedBook.Name,
		Author:      ReturnedBook.Author.FirstName + " " + ReturnedBook.Author.LastName,
		Category:    ReturnedBook.Category,
		Volume:      ReturnedBook.Volume,
		PublishedAt: ReturnedBook.PublishedAt,
		Summary:     ReturnedBook.Summary,
		Publisher:   ReturnedBook.Publisher,
	}

	// add related contents to the result
	for _, c := range *Contents {
		Book.TableOfContents = append(Book.TableOfContents, c.ContentName)
	}

	// marshal the data
	resBody, err := json.Marshal(Book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bm.Logger.WithError(err).Error("error trying to marshal the response message")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBody)

}

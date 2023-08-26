package handler

import (
	"bookman/authenticate"
	"encoding/json"
	"net/http"
)

func (bm *BookManagerServer) GetAllBooksHandler(w http.ResponseWriter, r *http.Request) {
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

	// check if the username exists
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

	// retrieve all books from database
	var books Books
	books.Books, err = bm.DB.GetAllBooks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bm.Logger.WithError(err).Error("error retrieving all the books from the database")
		return
	}

	// append all records to a  preferred-form of data structure
	var AllBooks []*GetAllBooksResponseBody
	for _, b := range *books.Books {
		temp := &GetAllBooksResponseBody{
			BookId:      b.ID,
			Name:        b.Name,
			Author:      b.Author.FirstName + " " + b.Author.LastName,
			Category:    b.Category,
			Volume:      b.Volume,
			PublishedAt: b.PublishedAt,
			Summary:     b.Summary,
			Publisher:   b.Publisher,
		}
		AllBooks = append(AllBooks, temp)
	}

	// marshalling response body
	resBody, err := json.Marshal(map[string]interface{}{
		"books": AllBooks,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bm.Logger.WithError(err).Error("error trying to marshal the response message")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

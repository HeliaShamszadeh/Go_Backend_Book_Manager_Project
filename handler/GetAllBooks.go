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
		bm.Logger.WithError(err).Warn("could not log in the user")
		response, _ := json.Marshal(map[string]interface{}{"message": err.Error()})
		w.Write(response)
		return
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
	var AllBooks []*GetBooksResponseBody
	for _, b := range *books.Books {
		temp := &GetBooksResponseBody{
			Name:        b.Name,
			Author:      b.FirstName + " " + b.LastName,
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
		bm.Logger.WithError(err).Warn("error writing response body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

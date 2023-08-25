package handler

import (
	"bookman/authenticate"
	"bookman/db"
	"encoding/json"
	"io"
	"net/http"
)

func (bm *BookManagerServer) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	// check the request method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get access token from header
	AuthorizationToken := r.Header.Get("Authorization")
	if AuthorizationToken == "" {
		resBody, _ := json.Marshal(map[string]interface{}{
			"message": authenticate.InvalidTokenErr,
		})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resBody)
		return
	}

	// get the username using the access token
	username, err := bm.Authenticate.GetUsernameByToken(AuthorizationToken)
	if err != nil {
		if err == authenticate.CannotValidateToken {
			bm.Logger.WithError(err).Warn()
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err == authenticate.InvalidTokenErr {
			bm.Logger.WithError(err).Warn()
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			bm.Logger.WithError(err).Warn("there was a problem with access token")
			return
		}
	}

	// get user account by access token
	account, err := bm.DB.GetUserByUsername(username)
	if err != nil {
		if err == db.UserNameNotFoundError {
			bm.Logger.WithError(err).Warn()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		bm.Logger.WithError(err).Warn("error while retrieving user from db")
		return
	}

	// check if request body is nil
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// read the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		bm.Logger.WithError(err).Warn("cannot read the request data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// unmarshalling request body data
	var NewBook CreateBookRequestBody
	err = json.Unmarshal(reqBody, &NewBook)
	if err != nil {
		bm.Logger.WithError(err).Warn("error while unmarshalling book")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// insert new book to database by calling GormDB insert function
	newBook := &db.Book{
		Name:        NewBook.Name,
		Category:    NewBook.Category,
		Volume:      NewBook.Volume,
		PublishedAt: NewBook.PublishedAt,
		Summary:     NewBook.Summary,
		Publisher:   NewBook.Publisher,
		FirstName:   NewBook.Author.FirstName,
		LastName:    NewBook.Author.LastName,
		Birthday:    NewBook.Author.Birthday,
		Nationality: NewBook.Author.Nationality,
		UserID:      account.ID,
	}

	err = bm.DB.CreateNewBook(newBook)
	if err != nil {
		bm.Logger.WithError(err).Warn("error while inserting new book")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(map[string]interface{}{
		"message": "book was created successfully",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bm.Logger.WithError(err).Error("error trying to marshal the response message")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

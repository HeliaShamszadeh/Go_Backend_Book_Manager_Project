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

	// get the username using the access token
	username, err := bm.Authenticate.GetUsernameByToken(AuthorizationToken)
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

	// get user account by access token
	account, err := bm.DB.GetUserByUsername(username)
	if err != nil {
		if err == db.UserNameNotFoundErr {
			bm.Logger.WithError(err).Warn()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		bm.Logger.WithError(err).Warn("error while retrieving user from db (get user by username")
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
	var NewBook db.Book
	var Table TableOfContents

	err = json.Unmarshal(reqBody, &Table)
	if err != nil {
		bm.Logger.WithError(err).Warn("error while unmarshalling table of contents")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &NewBook)
	if err != nil {
		bm.Logger.WithError(err).Warn("error while unmarshalling book")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	NewBook.UserID = account.ID

	// add each content to the book instance
	for _, content := range Table.Contents {
		NewBook.TableOfContents = append(NewBook.TableOfContents, db.Content{ContentName: content})
	}

	// insert new book to database
	err = bm.DB.CreateNewBook(&NewBook)
	if err != nil {
		bm.Logger.WithError(err).Warn("error while inserting new book (create new book)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create response body
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

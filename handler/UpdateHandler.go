package handler

import (
	"bookman/authenticate"
	"bookman/db"
	"encoding/json"
	"io"
	"net/http"
	"path"
	"strconv"
)

func (bm *BookManagerServer) UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	// check the request method
	if r.Method != http.MethodPut {
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

	// get user's account by username
	account, err := bm.DB.GetUserByUsername(username)
	if err != nil {
		bm.Logger.WithError(err).Warn("error finding account by username")
		w.WriteHeader(http.StatusInternalServerError) // ???
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
		if err == db.BookNotFoundErr {
			w.WriteHeader(http.StatusBadRequest)
			resBody, _ := json.Marshal(map[string]interface{}{"message": "book not found"})
			w.Write(resBody)
			return
		}
		bm.Logger.WithError(err).Warn("error reading from database (GetBookById)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if the request body is nil
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// read the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// unmarshalling data
	var urb UpdateRequestBody
	err = json.Unmarshal(reqBody, &urb)
	if err != nil {
		bm.Logger.WithError(err).Warn("error in unmarshalling request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the logged-in user owns this book
	if ReturnedBook.UserID != account.ID {
		w.WriteHeader(http.StatusForbidden)
		resBody, _ := json.Marshal(map[string]interface{}{"message": "access denied"})
		w.Write(resBody)
		return
	}

	// call GormDB update function
	err = bm.DB.UpdateBook(ReturnedBook, urb.Name, urb.Category)
	if err != nil {
		bm.Logger.WithError(err).Warn("error updating book (UpdateBookHandler)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// marshalling a success message
	response, err := json.Marshal(map[string]interface{}{
		"message": "book was updated successfully!"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bm.Logger.WithError(err).Error("error trying to marshal the response message")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

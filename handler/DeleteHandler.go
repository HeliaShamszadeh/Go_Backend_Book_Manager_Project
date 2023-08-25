package handler

import (
	"bookman/authenticate"
	"bookman/db"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
)

func (bm *BookManagerServer) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	// check the request method
	if r.Method != http.MethodDelete {
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
	username, err := bm.Authenticate.GetUsernameByToken(AuthorizationToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bm.Logger.WithError(err).Warn("could not login the user")
		return
	}
	// get user's account info by username
	account, err := bm.DB.GetUserByUsername(username)
	if err != nil {
		bm.Logger.WithError(err).Warn("error finding account by username")
		w.WriteHeader(http.StatusBadRequest) // ???
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
		w.WriteHeader(http.StatusBadRequest)
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
		bm.Logger.WithError(err).Warn("error reading from database (GetBookById)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if the logged-in user owns this book
	if ReturnedBook.UserID != account.ID {
		w.WriteHeader(http.StatusForbidden)
		resBody, _ := json.Marshal(map[string]interface{}{"message": "access denied"})
		w.Write(resBody)
		return
	}

	// calling GormDB delete function
	err = bm.DB.DeleteBookById(BookIdInt)
	if err != nil {
		bm.Logger.WithError(err).Warn("error deleting object")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// marshall response body
	response, err := json.Marshal(map[string]interface{}{"message": "book was deleted succesfully!"})
	if err != nil {
		bm.Logger.WithError(err).Warn("error marshaling response body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

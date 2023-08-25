package handler

import (
	"bookman/db"
	"encoding/json"
	"io"
	"net/http"
)

func (bm *BookManagerServer) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	// check the request method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// read the request body
	reqData, err := io.ReadAll(r.Body)
	if reqData == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bm.Logger.WithError(err).Warn("cannot read request body")
		return
	}

	// create new model for new user and unmarshalling data
	var NewUser db.User
	err = json.Unmarshal(reqData, &NewUser)
	if err != nil {
		bm.Logger.WithError(err).Warn("error unmarshalling request body (create new user)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call create new user GormDB function
	err = bm.DB.CreateNewUser(&NewUser)
	if err != nil {
		switch err {
		case db.DuplicateUsernameError:
			w.WriteHeader(http.StatusConflict)
			return
		case db.DuplicateEmailError:
			w.WriteHeader(http.StatusConflict)
			return
		case db.DuplicatePhoneNumberError:
			w.WriteHeader(http.StatusConflict)
			return
		default:
			w.WriteHeader(http.StatusBadRequest)
			bm.Logger.WithError(err).Warn("error adding new user to database")
			return

		}
	}

	// write response message
	response := map[string]interface{}{
		"message": "User has been created successfully!",
	}
	resBody, _ := json.Marshal(response)
	w.WriteHeader(http.StatusCreated)
	w.Write(resBody)
}

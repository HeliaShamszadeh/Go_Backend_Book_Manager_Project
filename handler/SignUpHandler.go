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

	// check if the request body is nil
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
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
		case db.DuplicateUsernameErr:
			w.WriteHeader(http.StatusConflict)
			response, _ := json.Marshal(map[string]interface{}{"message": db.DuplicateUsernameErr.Error()})
			w.Write(response)
			return
		case db.DuplicateEmailErr:
			w.WriteHeader(http.StatusConflict)
			response, _ := json.Marshal(map[string]interface{}{"message": db.DuplicateEmailErr.Error()})
			w.Write(response)
			return
		case db.DuplicatePhoneNumberErr:
			w.WriteHeader(http.StatusConflict)
			response, _ := json.Marshal(map[string]interface{}{"message": db.DuplicatePhoneNumberErr.Error()})
			w.Write(response)
			return
		case db.GenderNotAllowedErr:
			w.WriteHeader(http.StatusBadRequest)
			response, _ := json.Marshal(map[string]interface{}{"message": db.GenderNotAllowedErr.Error()})
			w.Write(response)
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

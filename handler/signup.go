package handler

import (
	"bookman/db"
	"encoding/json"
	"io"
	"net/http"
)

func (bm *BookManagerServer) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bm.Logger.WithError(err).Warn("Cannot Read Request Body")
		return
	}

	var NewUser db.User
	err = json.Unmarshal(reqData, &NewUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bm.Logger.WithError(err).Warn("Cannot Unmarshal Request Body")
		return
	}

	err = bm.DB.CreateNewUser(&NewUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response1 := map[string]interface{}{
			"message": err.Error(),
		}
		responseBody, _ := json.Marshal(response1)
		w.Write([]byte(responseBody))
		bm.Logger.WithError(err).Warn("Cannot Add New User to Database")
		return
	}

	response := map[string]interface{}{
		"message": "User has been created successfully!",
	}
	resBody, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

package handler

import (
	"bookman/db"
	"encoding/json"
	"io"
	"net/http"
)

type signupRequestBody struct {
	Username    string `json:"user_name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
}

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

	var srb signupRequestBody
	err = json.Unmarshal(reqData, &srb)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bm.Logger.WithError(err).Warn("Cannot Unmarshal Request Body")
		return
	}

	NewUser := &db.User{
		Username:    srb.Username,
		Email:       srb.Email,
		Password:    srb.Password,
		FirstName:   srb.FirstName,
		LastName:    srb.LastName,
		PhoneNumber: srb.PhoneNumber,
		Gender:      srb.Gender,
	}

	err = bm.DB.CreateNewUser(NewUser)
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

package handler

import (
	"bookman/authenticate"
	"encoding/json"
	"io"
	"net/http"
)

func (bm *BookManagerServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// check the request method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// check if request body is nil
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// read the request body
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		bm.Logger.WithError(err).Warn("cannot read the request data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// unmarshall request body data
	var lrb LoginRequestBody
	err = json.Unmarshal(reqData, &lrb)
	if err != nil {
		bm.Logger.WithError(err).Warn("cannot unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// login authentication
	token, err := bm.Authenticate.Login(authenticate.Credential{
		Username: lrb.Username,
		Password: lrb.Password,
	})
	if err != nil {
		if err == authenticate.IncorrectPasswordErr {
			w.WriteHeader(http.StatusForbidden)
			response, _ := json.Marshal(map[string]interface{}{"message": "incorrect password"})
			w.Write(response)
			return
		} else {
			w.WriteHeader(http.StatusForbidden)
			bm.Logger.WithError(err).Warn("error logging the user in")
			response, _ := json.Marshal(map[string]interface{}{"message": err.Error()})
			w.Write(response)
			return
		}
	}

	// write response body
	resBody, _ := json.Marshal(map[string]interface{}{
		"access_token": token.TokenString,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

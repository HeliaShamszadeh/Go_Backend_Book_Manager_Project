package handler

import (
	"bookman/authenticate"
	"encoding/json"
	"io"
	"net/http"
)

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (bm *BookManagerServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		bm.Logger.Warn("request body is empty")
		return
	}

	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		bm.Logger.WithError(err).Warn("cannot read the request data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var lrb loginRequestBody
	err = json.Unmarshal(reqData, &lrb)
	if err != nil {
		bm.Logger.WithError(err).Warn("cannot unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := bm.Authenticate.Login(authenticate.Credential{
		Username: lrb.Username,
		Password: lrb.Password,
	})

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	response := map[string]interface{}{
		"access_token": token.TokenString,
	}

	resBody, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

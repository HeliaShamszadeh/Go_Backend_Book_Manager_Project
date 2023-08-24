package handler

import (
	"bookman/authenticate"
	"bookman/db"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type createBookRequestBody struct {
	Name            string    `json:"name"`
	Author          author    `json:"author"`
	Category        string    `json:"category"`
	Volume          int       `json:"volume"`
	PublishedAt     time.Time `json:"published_at"`
	Summary         string    `json:"summary"`
	TableOfContents []string  `json:"table_of_contents"`
	Publisher       string    `json:"publisher"`
}

type author struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Birthday    time.Time `json:"birthday"`
	Nationality string    `json:"nationality"`
}

func (bm *BookManagerServer) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	username, err := bm.Authenticate.GetUsernameByToken(AuthorizationToken)
	if err != nil {
		if err == authenticate.CannotValidateToken {
			bm.Logger.WithError(err).Warn()
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			bm.Logger.WithError(err).Warn("there was a problem with access token")
			return
		}
	}

	account, err := bm.DB.GetUserByUsername(username)
	if err == db.UserNameNotFoundError {
		bm.Logger.WithError(err).Warn()
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bm.Logger.WithError(err).Warn("error while retrieving user from db")
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var NewBook createBookRequestBody
	err = json.Unmarshal(reqBody, &NewBook)
	if err != nil {
		bm.Logger.WithError(err).Warn("error while unmarshalling book")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	message := map[string]interface{}{
		"message": "book was created successfully",
	}

	respone, err := json.Marshal(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bm.Logger.WithError(err).Error("error trying to marshal the respone message")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respone)
}

package handler

import (
	"bookman/authenticate"
	"bookman/db"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type responseBody struct {
	Books []book
}

type book struct {
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Category    string    `json:"category"`
	Volume      int       `json:"volume"`
	PublishedAt time.Time `json:"published_at"`
	Summary     string    `json:"summary"`
	Publisher   string    `json:"publisher"`
}

type row struct {
	name         string
	first_name   string
	last_name    string
	category     string
	volume       int
	published_at time.Time
	summary      string
	publisher    string
}

func (bm *BookManagerServer) GetAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	var books []book
	//rows, err := bm.DB.Db.Table("books").Select("name, first_name, last_name, category, volume, published_at, summary, publisher").Rows()
	rows, err := bm.DB.Db.Model(&db.Book{}).Select("name, first_name, last_name, category, volume, published_at, summary, publisher").Rows()
	fmt.Println(rows)
	if err != nil {
		bm.Logger.WithError(err).Warn("error retrieving objects from db")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var allbooks responseBody
	for rows.Next() {
		var (
			name         string
			first_name   string
			last_name    string
			category     string
			volume       int
			published_at time.Time
			summary      string
			publisher    string
		)

		err = rows.Scan(&name, &first_name, &last_name, &category, &volume, &published_at, &summary, &publisher)
		if err != nil {
			bm.Logger.WithError(err).Warn("error while scanning rows")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		book := book{
			Name:        name,
			Author:      first_name + " " + last_name,
			Category:    category,
			Volume:      volume,
			PublishedAt: published_at,
			Summary:     summary,
			Publisher:   publisher,
		}
		_ = append(allbooks.Books, book)
	}

	resBody, err := json.Marshal(books)
	if err != nil {
		bm.Logger.WithError(err).Warn("error writing response body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

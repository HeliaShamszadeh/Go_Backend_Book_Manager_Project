package main

import (
	"bookman/authenticate"
	"bookman/config"
	"bookman/db"
	"bookman/handler"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err.Error())
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	gormDB, err := db.NewGormDB(cfg)
	if err != nil {
		logger.WithError(err).Fatalln("ERROR in the Database Migration")
	}
	logger.Infoln("Migrated Tables and Models Successfully!")
	err1, err2 := gormDB.CreateSchemas()
	if err1 != nil {
		logger.WithError(err1).Fatalln("ERROR Occurred While Creating User Model")
	}
	if err2 != nil {
		logger.WithError(err2).Fatalln("ERROR Occurred While Creating Book Model")
	}

	auth, err := authenticate.NewAuth(gormDB, time.Minute*10, logger)
	bookManagerServer := handler.BookManagerServer{DB: gormDB, Logger: logger, Authenticate: auth}

	http.HandleFunc("/api/v1/auth/signup", bookManagerServer.SignUpHandler)
	http.HandleFunc("/api/v1/auth/login", bookManagerServer.LoginHandler)
	http.HandleFunc("/api/v1/createbooks", bookManagerServer.CreateBookHandler)
	http.HandleFunc("/api/v1/books", bookManagerServer.GetAllBooksHandler)

	logger.WithError(http.ListenAndServe(":8080", nil)).Fatalln("can not setup the server")
}

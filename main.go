package main

import (
	"bookman/config"
	"bookman/db"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
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
}

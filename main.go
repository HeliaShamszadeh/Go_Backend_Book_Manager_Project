package main

import (
	"bookman/config"
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

}

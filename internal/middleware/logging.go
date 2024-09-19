package middleware

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func InitLogger() {
	log.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("internal/logs/content-app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		Logger.Println("Failed to log to file, using default stderr")
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(file)
	}
	log.SetLevel(logrus.InfoLevel)

}
func GetLogger() *logrus.Logger {
	return log
}

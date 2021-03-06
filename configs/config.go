package config

import (
	"os"
	"sync"

	"github.com/furqonzt99/snackbox/constants"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/xendit/xendit-go"
)

type AppConfig struct {
	Port     string
	Database struct {
		Driver   string
		Name     string
		Host     string
		Port     string
		Username string
		Password string
	}
}

var lock = &sync.Mutex{}
var appConfig *AppConfig
var Mode string

func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig()
	}

	return appConfig
}

func initConfig() *AppConfig {

	err := godotenv.Load()

	if err != nil {
		log.Info("Error loading .env file")
	}

	var defaultConfig AppConfig
	defaultConfig.Port = os.Getenv("APP_PORT")
	defaultConfig.Database.Driver = os.Getenv("DB_DRIVER")
	defaultConfig.Database.Name = os.Getenv("DB_NAME")
	defaultConfig.Database.Host = os.Getenv("DB_HOST")
	defaultConfig.Database.Port = os.Getenv("DB_PORT")
	defaultConfig.Database.Username = os.Getenv("DB_USERNAME")
	defaultConfig.Database.Password = os.Getenv("DB_PASSWORD")

	constants.JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	constants.XENDIT_CALLBACK_TOKEN = os.Getenv("XENDIT_CALLBACK_TOKEN")

	constants.AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
	constants.AWS_ACCESS_SECRET_KEY = os.Getenv("AWS_ACCESS_SECRET_KEY")
	constants.S3_REGION = os.Getenv("S3_REGION")
	constants.S3_BUCKET = os.Getenv("S3_BUCKET")
	constants.LINK_TEMPLATE = os.Getenv("LINK_TEMPLATE")

	xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

	Mode = os.Getenv("MODE")

	return &defaultConfig
}

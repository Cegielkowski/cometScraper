package config

import (
	"cometScraper/tools/scraper/pkg/element"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPORT     string
	DatabaseURL    string
	CacheURL       string
	LoggerLevel    string
	ContextTimeout int
	Elements       element.Elements
}

// LoadConfig will load config from environment variable
func LoadConfig() (config *Config) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	serverPORT := os.Getenv("SERVER_PORT")
	databaseURL := os.Getenv("DATABASE_URL")
	cacheURL := os.Getenv("CACHE_URL")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	elementsInputPath := os.Getenv("ELEMENTS_INPUT_PATH")
	contextTimeout, _ := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))

	fileContent, err := os.Open(elementsInputPath)

	if err != nil {
		log.Println(err)
		return
	}

	defer fileContent.Close()

	elements, err := element.NewElement(fileContent)
	if err != nil {
		panic(err)
	}
	return &Config{
		ServerPORT:     serverPORT,
		DatabaseURL:    databaseURL,
		CacheURL:       cacheURL,
		LoggerLevel:    loggerLevel,
		ContextTimeout: contextTimeout,
		Elements:       elements,
	}
}

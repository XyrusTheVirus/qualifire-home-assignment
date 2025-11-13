package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// getConfigPath returns the path to the config files
func getConfigPath() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	isTest := EnvInt("IS_TEST", "0") == 1
	var configPath string
	if isTest {
		configPath = wd + "/../internal/configs"
	} else {
		configPath = wd + "/internal/configs"
	}
	return configPath
}

// LoadConfig loadConfig loads the config file
func LoadConfig() {
	viper.SetConfigName("keys")
	viper.SetConfigType("json")
	viper.AddConfigPath(getConfigPath())
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

// Config returns the value of the key in the config file
func Config(key string, fallback string) any {
	value := viper.Get(key)
	if value == "" {
		return fallback
	}
	return value

}

// LoadEnv loadEnv loads the .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading/ .env file")
	}
}

// Env returns the value of the key in the .env file
func Env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// EnvInt returns the value of the key in the .env file as an int
func EnvInt(key string, fallback string) int {
	value, _ := strconv.Atoi(Env(key, fallback))
	return value
}

// IsDevelopment returns true if the app is in development mode
func IsDevelopment() bool {
	return Env("APP_ENV", "") == "dev"
}

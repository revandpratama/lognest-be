package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	APP_ENV string `mapstructure:"APP_ENV"`

	REST_PORT string `mapstructure:"REST_PORT"`

	AUTH4ME_URL string `mapstructure:"AUTH4ME_URL"`

	CORS_ALLOWED_ORIGINS string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CORS_ALLOWED_HEADERS string `mapstructure:"CORS_ALLOWED_HEADERS"`

	COOKIE_SECRET string `mapstructure:"COOKIE_SECRET"`

	JWT_SECRET            string `mapstructure:"JWT_SECRET"`
	JWT_EXPIRATION_SECOND string `mapstructure:"JWT_EXPIRATION_SECOND"`

	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     string `mapstructure:"DB_PORT"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_NAME     string `mapstructure:"DB_NAME"`

	LOGNEST_SCHEMA string `mapstructure:"LOGNEST_SCHEMA"`
	AUTH4ME_SCHEMA string `mapstructure:"AUTH4ME_SCHEMA"`

	AZURE_STORAGE_CONNECTION_STRING              string `mapstructure:"AZURE_STORAGE_CONNECTION_STRING"`
	AZURE_STORAGE_CONTAINER_NAME                 string `mapstructure:"AZURE_STORAGE_CONTAINER_NAME"`
	AZURE_STORAGE_URL_EXPIRY_DURATION_IN_MINUTES string `mapstructure:"AZURE_STORAGE_URL_EXPIRY_DURATION_IN_MINUTES"`
}

var ENV Config

func LoadConfig() error {

	// * For local development
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetDefault("REST_PORT", "8080")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("No .env file found, relying on runtime environment variables.")
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// * For docker environment
	v := reflect.ValueOf(ENV)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// Get the struct tag, which is the name of the environment variable
		key := t.Field(i).Tag.Get("mapstructure")
		if key != "" {
			// Bind the key to the corresponding environment variable
			// e.g., tells Viper that the key "DB_HOST" should be read from the env var "DB_HOST"
			if err := viper.BindEnv(key); err != nil {
				return err
			}
		}
	}
	viper.AutomaticEnv()

	// Unmarshal into the ENV variable
	if err := viper.Unmarshal(&ENV); err != nil {
		return err
	}

	return nil
}

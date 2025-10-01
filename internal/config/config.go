package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Port           int
	AppEnv         string // development, production, qa
	ImageUploadDir string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Debug().Msg("error loading .env file")
	}

	viper.AutomaticEnv()

	port := mustGetInt("PORT")
	appEnv := mustGetString("APP_ENV")

	return &Config{
		Port:           port,
		AppEnv:         appEnv,
		ImageUploadDir: "./uploads",
	}

}

func mustGetInt(key string) int {
	v := viper.GetInt(key)
	if v == 0 {
		log.Fatal().Msgf("The env var %s must be set.", key)
	}

	return v
}

func mustGetString(key string) string {
	v := viper.GetString(key)
	if v == "" {
		log.Fatal().Msgf("The env var %s must be set.", key)
	}

	return v
}

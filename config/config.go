package config

import (
	"log"

	"github.com/spf13/viper"
)

type EnvModel struct {
	Port       string `mapstructure:"PORT"`
	DbName     string `mapstructure:"DB_NAME"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     string `mapstructure:"DB_PORT"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
}

func InitConfig() (configs *EnvModel) {

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error loading env env variables", err)
	}

	if err := viper.Unmarshal(&configs); err != nil {
		log.Fatal("Error while unmarshalling loaded variables into struct")
	}

	return
}

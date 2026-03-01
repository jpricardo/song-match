package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	viper.Set("APP_ENV", viper.Get("APP_ENV"))
	viper.Set("SERVER_ADDRESS", viper.Get("SERVER_ADDRESS"))
	viper.Set("CONTEXT_TIMEOUT", viper.Get("CONTEXT_TIMEOUT"))
	viper.Set("DB_HOST", viper.Get("DB_HOST"))
	viper.Set("DB_PORT", viper.Get("DB_PORT"))
	viper.Set("DB_USER", viper.Get("DB_USER"))
	viper.Set("DB_PASS", viper.Get("DB_PASS"))
	viper.Set("DB_NAME", viper.Get("DB_NAME"))
	viper.Set("ACCESS_TOKEN_EXPIRY_HOUR", viper.Get("ACCESS_TOKEN_EXPIRY_HOUR"))
	viper.Set("REFRESH_TOKEN_EXPIRY_HOUR", viper.Get("REFRESH_TOKEN_EXPIRY_HOUR"))
	viper.Set("ACCESS_TOKEN_SECRET", viper.Get("ACCESS_TOKEN_SECRET"))
	viper.Set("REFRESH_TOKEN_SECRET", viper.Get("REFRESH_TOKEN_SECRET"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Can't find the env file: ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}

package config

import (
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

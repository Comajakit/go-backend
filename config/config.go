package config

import (
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yamml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

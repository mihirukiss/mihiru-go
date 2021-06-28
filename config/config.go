package config

import (
	"github.com/spf13/viper"
	"log"
)

var configs *viper.Viper

func Init(env string) {
	var err error
	configs = viper.New()
	configs.SetConfigType("yaml")
	configs.SetConfigName(env)
	configs.AddConfigPath("../config/")
	configs.AddConfigPath("config/")
	err = configs.ReadInConfig()
	if err != nil {
		log.Fatal("Read config file fail", err.Error())
	}
}

func GetConfigs() *viper.Viper {
	return configs
}

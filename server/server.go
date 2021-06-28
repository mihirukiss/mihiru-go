package server

import (
	"log"
	"mihiru-go/config"
	"mihiru-go/database"
)

func Init() {
	configs := config.GetConfigs()
	mongoDatabase, err := database.New(configs.GetString("database.uri"), configs.GetString("database.username"), configs.GetString("database.password"), configs.GetString("database.dbname"))
	if err != nil {
		log.Fatal(err.Error())
	}
	r := NewRouter(mongoDatabase)
	err = r.Run(configs.GetStringSlice("server.addr")...)
	if err != nil {
		log.Fatal(err.Error())
	}
}

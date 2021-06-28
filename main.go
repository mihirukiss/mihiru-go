package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"mihiru-go/config"
	"mihiru-go/server"
	"os"
)

func main() {
	env := flag.String("e", "dev", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*env)
	gin.SetMode(config.GetConfigs().GetString("gin.mode"))
	server.Init()
}

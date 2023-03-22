package main

import (
	"log"
	"web_adb/routes"

	"github.com/gin-gonic/gin"
)

func checkError(e error, args ...any) {
	if e != nil {
		args = append(args, e)
		log.Fatal(args...)
	}
}

func main() {
	var server = gin.New()
	server.Use(gin.Recovery(), gin.Logger())
	routes.Home(server)
	routes.Files(server)
	server.Run(":http")
}

package main

import (
	"pokemon/config"
	db "pokemon/database"
	"pokemon/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	db.InitDB()
	r := gin.Default()
	r.POST("/users", handlers.RegisterHandler)
	r.POST("/del-users", handlers.DeleteRecentUserHandlerGin)

	r.Run(":3000")

}

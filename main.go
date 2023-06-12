package main

import (
	"log"
	"net/http"
	"pokemon/config"
	db "pokemon/database"
	"pokemon/routers"
)

func main() {
	config.InitConfig()
	db.InitDB()

	router := routers.SetupRouter()

	log.Fatal(http.ListenAndServe(":3000", router))

	// r := gin.Default()
	// config.InitConfig()
	// text := viper.Get("DB_HOST")

	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": text,
	// 	})
	// })

	// r.Run(":3000") // Start the Gin server
}

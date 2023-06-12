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

}

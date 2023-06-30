package main

import (
	"go-backend/config"
	db "go-backend/database"
	"go-backend/routers"
)

func main() {
	config.InitConfig()
	db.InitDB()
	r := routers.SetupRouter()
	r.Run(":3000")

}

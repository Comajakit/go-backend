package main

import (
	"go-backend/config"
	db "go-backend/database"
	"go-backend/routers"

	"github.com/gin-contrib/sessions/cookie"
)

func main() {
	config.InitConfig()
	db.InitDB()
	// sessionStore := config.GetSessionStore()
	store := cookie.NewStore([]byte("secret"))

	r := routers.SetupRouter(store)

	r.Run(":3000")

}

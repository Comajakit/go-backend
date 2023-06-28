package main

import (
	"go-backend/config"
	db "go-backend/database"
	"go-backend/handlers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	db.InitDB()
	r := gin.Default()
	store := cookie.NewStore([]byte("secret")) // Replace "secret" with your own secret key
	r.Use(sessions.Sessions("mysession", store))
	r.POST("/users", handlers.RegisterHandler)
	r.POST("/del-users", handlers.DeleteRecentUserHandlerGin)
	r.POST("/login", handlers.UserLogin)
	r.GET("/protected", handlers.ProtectedRoute)

	r.Run(":3000")

}

package routers

import (
	"go-backend/handlers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret")) // Replace "secret" with your own secret key
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/users", handlers.RegisterHandler)
	r.POST("/del-users", handlers.DeleteRecentUserHandlerGin)
	r.POST("/login", handlers.UserLogin)
	r.GET("/protected", handlers.ProtectedRoute)

	return r
}

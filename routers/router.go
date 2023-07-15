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

	user_prefix_path := "v1/user"
	port_prefix_path := "v1/port"
	//user
	r.POST(user_prefix_path+"/create", handlers.RegisterHandler)
	r.POST(user_prefix_path+"/delete", handlers.DeleteRecentUserHandlerGin)
	r.POST(user_prefix_path+"/login", handlers.UserLogin)

	//trip
	r.GET("/protected", handlers.ProtectedRoute)

	//port
	r.POST(port_prefix_path+"/create", handlers.CreatePort)
	r.POST(port_prefix_path+"/add-strategy", handlers.AddPortStrategy)
	r.POST(port_prefix_path+"/add-stock", handlers.AddStock)
	r.POST(port_prefix_path+"/check", handlers.CheckToken)

	return r
}

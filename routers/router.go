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
	trip_prefix_path := "v1/trip"
	//user
	r.POST(user_prefix_path+"/create", handlers.RegisterHandler)
	r.POST(user_prefix_path+"/delete", handlers.DeleteRecentUserHandlerGin)
	r.POST(user_prefix_path+"/login", handlers.UserLogin)

	r.POST(trip_prefix_path+"/create", handlers.CreateTrip)
	r.GET("/protected", handlers.ProtectedRoute)

	r.POST("/setvalue", handlers.SetValueHandler)
	r.GET("/getvalue", handlers.GetValueHandler)

	return r
}

package routers

import (
	"go-backend/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetupRouter(store sessions.Store) *gin.Engine {
	r := gin.Default()
	// Use the CORS middleware with Gin
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"*", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	r.Use(sessions.Sessions("testsession", store))
	user_prefix_path := "v1/user"
	port_prefix_path := "v1/port"
	//user
	r.POST(user_prefix_path+"/create", handlers.RegisterHandler)
	r.DELETE(user_prefix_path+"/delete", handlers.DeleteRecentUserHandlerGin)
	r.POST(user_prefix_path+"/login", handlers.UserLogin)
	r.GET(user_prefix_path+"/username", handlers.NameFromToken)

	//trip
	r.GET("/protected", handlers.ProtectedRoute)

	//port
	r.GET(port_prefix_path+"/get-dividend", handlers.GetCurrentDivPercent)
	r.GET(port_prefix_path+"/get-port", handlers.GetStock)

	r.POST(port_prefix_path+"/create", handlers.CreatePort)
	r.POST(port_prefix_path+"/add-strategy", handlers.CreatePortStrategy)
	r.POST(port_prefix_path+"/add-stock", handlers.AddStock)
	r.POST(port_prefix_path+"/summary", handlers.SummaryPort)

	r.PUT(port_prefix_path+"/update-stock", handlers.UpdateStock)
	r.PUT(port_prefix_path+"/update-strategy", handlers.UpdateStrategy)

	r.POST(port_prefix_path+"/delete-stock", handlers.DeleteStock)

	return r
}

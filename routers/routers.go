package routers

import (
	"pokemon/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
	router.HandleFunc("/del-users", handlers.DeleteRecentUserHandler).Methods("POST")

	return router
}

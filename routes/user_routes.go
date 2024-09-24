package routes

import (
	"github.com/chamanbetra/user-management-app/controllers"
	"github.com/gorilla/mux"
)

func router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users", controllers.CreateUser).Methods("POST")

	r.HandleFunc("/user", controllers.BasicAuth(controllers.GetUser)).Methods("GET")
	r.HandleFunc("/user", controllers.BasicAuth(controllers.UpdateUser)).Methods("PUT")
	r.HandleFunc("/user", controllers.BasicAuth(controllers.DeleteUser)).Methods("DELETE")

	return r
}

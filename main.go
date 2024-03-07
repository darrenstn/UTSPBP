package main

import (
	"UTS/controllers"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	//loadEnv()

	router := mux.NewRouter()
	router.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/v2/user", controllers.InsertUserV2).Methods("POST")
	router.HandleFunc("/v1/user", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/v2/user", controllers.UpdateUserV2).Methods("PUT")
	router.HandleFunc("/user/{user_id}", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}

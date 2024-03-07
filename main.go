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
	router.HandleFunc("/v1/products", controllers.GetAllProducts).Methods("GET")
	router.HandleFunc("/v2/products", controllers.GetAllProductsV2).Methods("GET")
	router.HandleFunc("/transactions", controllers.GetAllTransactions).Methods("GET")
	router.HandleFunc("/v1/detail/transactions", controllers.GetAllDetailTransactions).Methods("GET")
	router.HandleFunc("/v2/detail/transactions", controllers.GetAllDetailTransactionsV2).Methods("GET")
	router.HandleFunc("/v1/user", controllers.InsertUser).Methods("POST")
	router.HandleFunc("/v2/user", controllers.InsertUserV2).Methods("POST")
	router.HandleFunc("/product", controllers.InsertProduct).Methods("POST")
	router.HandleFunc("/transaction", controllers.InsertTransaction).Methods("POST")
	router.HandleFunc("/v1/user", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/v2/user", controllers.UpdateUserV2).Methods("PUT")
	router.HandleFunc("/product", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/transaction", controllers.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/user/{user_id}", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/product/{product_id}", controllers.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/v1/transaction/{transaction_id}", controllers.DeleteTransaction).Methods("DELETE")
	router.HandleFunc("/v2/transaction/{transaction_id}", controllers.DeleteTransactionV2).Methods("DELETE")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}

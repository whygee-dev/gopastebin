package main

import (
	"gopastebin/db"
	"gopastebin/handler"
	"gopastebin/middleware"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)


func main() {
	db := db.CreateDb()

    router := mux.NewRouter()

	router.HandleFunc("/user/signup", handler.Signup(db)).Methods("POST")
	router.HandleFunc("/user/login", handler.Login(db)).Methods("POST")

	
    log.Fatal(http.ListenAndServe(":3333", middleware.AuthMiddleware(middleware.JsonContentTypeMiddleware(router))))
}




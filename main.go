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

	router.HandleFunc("/paste/create", handler.CreatePasteHandler(db)).Methods("PUT")
    router.HandleFunc("/paste/{id}", handler.GetPasteHandler(db)).Methods("GET")
	router.HandleFunc("/stats", handler.GetStatsHandler(db)).Methods("GET")

    log.Fatal(http.ListenAndServe(":3333", middleware.AuthMiddleware(middleware.JsonContentTypeMiddleware(router))))
}




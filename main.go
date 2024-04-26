package main

import (
	"gopastebin/consts"
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

	handler.SetupUserRoutes(db, router)
	handler.SetupPasteRoutes(db, router)

	log.Println("Server starting on port " + consts.GetPort())

    log.Fatal(http.ListenAndServe(":" + consts.GetPort(), middleware.AuthMiddleware(middleware.JsonContentTypeMiddleware(router))))
}




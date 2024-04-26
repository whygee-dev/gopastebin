package handler

import (
	"database/sql"
	"encoding/json"
	"gopastebin/models"
	"gopastebin/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupPasteRoutes(db *sql.DB, router *mux.Router) {
	router.HandleFunc("/paste/create", CreatePasteHandler(db)).Methods("PUT")
	router.HandleFunc("/paste/{id}", GetPasteHandler(db)).Methods("GET")
	router.HandleFunc("/stats", GetPasteStatsHandler(db)).Methods("GET")
}

func GetPasteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetPasteHandler")
		vars := mux.Vars(r)
		id := vars["id"]

		paste, err := service.GetPaste(db, id)

		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(&paste)
	}
}

func CreatePasteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("CreatePasteHandler")
		var body models.CreatePaste
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)

		if err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		id, err := service.CreatePaste(db, body)

		if err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

func GetPasteStatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetPasteStatsHandler")

		var stats models.Stats

		stats, err := service.GetStats(db)

		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		if err != nil {
			log.Println(err)

			
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(&stats)
	}
}

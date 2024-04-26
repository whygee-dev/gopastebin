package handler

import (
	"database/sql"
	"encoding/json"
	"gopastebin/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sqids/sqids-go"
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

		row := db.QueryRow("SELECT content, short_id, click_count, created_at FROM paste WHERE short_id = ?", id)

		var paste models.Paste

		err := row.Scan(&paste.Content, &paste.ShortID, &paste.ClickCount, &paste.CreatedAt)

		if err == sql.ErrNoRows {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		_, err = db.Exec("UPDATE paste SET click_count = click_count + 1 WHERE short_id = ?", id)

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

		row := db.QueryRow("SELECT MAX(id) as id FROM paste")

		var current_id int
		row.Scan(&current_id)

		s, _ := sqids.New()
		id, _ := s.Encode([]uint64{uint64(current_id + 1)})

		_, err = db.Exec("INSERT INTO paste (content, short_id, click_count) VALUES (?, ?, ?)", body.Content, id, 0)

		if err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

func GetPasteStatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetPasteStatsHandler")

		var stats models.Stats

		row := db.QueryRow("SELECT COUNT(*) as paste_count, SUM(click_count) as click_count, ROUND(AVG(click_count), 0) as avg_click_count FROM paste")

		err := row.Scan(&stats.TotalPastes, &stats.TotalClicks, &stats.AvgClicks)

		if err == sql.ErrNoRows {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(&stats)
	}
}

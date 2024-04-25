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

func GetPasteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetPasteHandler")

		vars := mux.Vars(r)
		id := vars["id"]

		row := db.QueryRow("SELECT content, short_id, click_count, created_at FROM paste WHERE short_id = ?", id)

		var paste models.Paste = models.Paste{}

		row.Scan(&paste.Content, &paste.ShortID, &paste.ClickCount, &paste.CreatedAt)

		_, err := db.Exec("UPDATE paste SET click_count = click_count + 1 WHERE short_id = ?", id)

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

		var body = models.CreatePaste{}

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

func GetStatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetStatsHandler")

		var stats = models.Stats{}

		row := db.QueryRow("SELECT COUNT(*) as paste_count, SUM(click_count) as click_count FROM paste")

		row.Scan(&stats.TotalPastes, &stats.TotalClicks)
		
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(&stats)
	}
}

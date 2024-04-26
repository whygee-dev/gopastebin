package service

import (
	"database/sql"
	"gopastebin/models"
	"log"

	"github.com/sqids/sqids-go"
)

func CreatePaste(db *sql.DB, data models.CreatePaste) (string, error) {
	row := db.QueryRow("SELECT MAX(id) as id FROM paste")

	var current_id int
	row.Scan(&current_id)

	s, _ := sqids.New()
	id, _ := s.Encode([]uint64{uint64(current_id + 1)})

	_, err := db.Exec("INSERT INTO paste (content, short_id, click_count) VALUES (?, ?, ?)", data.Content, id, 0)

	if err != nil {
		log.Println(err)

		return "", err
	}

	return id, nil
}

func GetPaste(db *sql.DB, id string) (models.Paste, error) {
	row := db.QueryRow("SELECT content, short_id, click_count, created_at FROM paste WHERE short_id = ?", id)
	var paste models.Paste
	err := row.Scan(&paste.Content, &paste.ShortID, &paste.ClickCount, &paste.CreatedAt)

	if err == sql.ErrNoRows {
		return models.Paste{}, err
	}

	_, err = db.Exec("UPDATE paste SET click_count = click_count + 1 WHERE short_id = ?", id)

	if err != nil {
		return models.Paste{}, err
	}

	return paste, nil
}
	

func GetStats(db *sql.DB) (models.Stats, error) {
	var stats models.Stats

	row := db.QueryRow("SELECT COUNT(*) as paste_count, SUM(click_count) as click_count, ROUND(AVG(click_count), 0) as avg_click_count FROM paste")

	err := row.Scan(&stats.TotalPastes, &stats.TotalClicks, &stats.AvgClicks)

	if err == sql.ErrNoRows {
		return models.Stats{}, err
	}

	return stats, nil
}
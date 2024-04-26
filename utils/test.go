package utils

import (
	"database/sql"
	"os"
)

func BuildUrl(path string, port string) string {
	return UrlFor("http", port, path)
}

func UrlFor(scheme string, serverPort string, path string) string {
	return scheme + "://localhost:" + serverPort + path
}

func TestTearDown(db *sql.DB) {
	wd, _ := os.Getwd()

	os.Remove(wd +  "/gopastebin.db")
}
package handler

import (
	"bytes"
	"encoding/json"
	"gopastebin/consts"
	"gopastebin/db"
	"gopastebin/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)


func TestGetPasteHandlerHappyPath(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupPasteRoutes(db, router)
	defer utils.TestTearDown(db)

	_, err := db.Exec("INSERT INTO paste (content, short_id, click_count) VALUES (?, ?, ?)", "content", "short_id", 1)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", utils.BuildUrl("/paste/short_id", consts.GetPort()), nil)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Response code was %v; want 200", res.Code)
	}

	var result map[string]interface{}
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&result)

	if result["content"] != "content" {
		t.Errorf("Response body was %v; want 'content'", result["content"])
	}

	if result["shortId"] != "short_id" {
		t.Errorf("Response body was %v; want 'short_id'", result["shortId"])
	}

	if int(result["clickCount"].(float64)) != 1 {
		t.Errorf("Response body was %v; want 1", result["clickCount"])
	}

	if result["createdAt"] == "" {
		t.Errorf("Response body was %v; want not empty", result["createdAt"])
	}

	paste_row_after := db.QueryRow("SELECT click_count FROM paste WHERE short_id = ?", "short_id")

	var click_count int

	paste_row_after.Scan(&click_count)

	if click_count != 2 {
		t.Errorf("Click count was %v; want 2", click_count)
	}
}

func TestGetPasteHandlerNotFound(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupPasteRoutes(db, router)
	defer utils.TestTearDown(db)

	_, err := db.Exec("INSERT INTO paste (content, short_id, click_count, expiry) VALUES (?, ?, ?, ?)", "content", "short_id", 1, "2020-01-01T00:00:00Z")

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", utils.BuildUrl("/paste/short_id", consts.GetPort()), nil)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("Response code was %v; want 404", res.Code)
	}
}

func TestGetPasteHandlerExpired(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupPasteRoutes(db, router)
	defer utils.TestTearDown(db)

	req, err := http.NewRequest("GET", utils.BuildUrl("/paste/short_id", consts.GetPort()), nil)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("Response code was %v; want 404", res.Code)
	}
}

func TestGetPasteInternalServerError(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupPasteRoutes(db, router)
	defer utils.TestTearDown(db)

	_, err := db.Exec(`DROP TABLE paste`)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", utils.BuildUrl("/paste/short_id", consts.GetPort()), nil)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("Response code was %v; want 500", res.Code)
	}
}

func TestCreatePasteHandlerHappyPath(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupPasteRoutes(db, router)
	defer utils.TestTearDown(db)

	body := map[string]string{
		"content": "content",
		"expiry": "2020-01-01T00:00:00Z",
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(body)

	req, err := http.NewRequest("PUT", utils.BuildUrl("/paste/create", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("Response code was %v; want 201", res.Code)
	}

	var result map[string]interface{}

	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&result)

	if result["id"] == "" {
		t.Errorf("Response body was %v; want not empty", result["shortId"])
	}

	paste_row := db.QueryRow("SELECT content, expiry FROM paste WHERE short_id = ?", result["id"])

	var content string
	var expiry string

	paste_row.Scan(&content, &expiry)

	if content != "content" {
		t.Errorf("Content was %v; want 'content'", content)
	}

	if expiry != "2020-01-01T00:00:00Z" {
		t.Errorf("Expiry was %v; want '2020-01-01T00:00:00Z'", expiry)
	}
}

func TestCreatePasteHandlerInvalidBody(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupPasteRoutes(db, router)
	defer utils.TestTearDown(db)

	buffer := new(bytes.Buffer)

	req, err := http.NewRequest("PUT", utils.BuildUrl("/paste/create", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Response code was %v; want 400", res.Code)
	}
}


func TestStatsHappyPath(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupPasteRoutes(db, router)
	defer utils.TestTearDown(db)

	_, err := db.Exec(`
		INSERT INTO paste (content, short_id, click_count) 
		VALUES 
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?)
	`, "content1", "short_id", 1, "content2", "short_id2", 2, "content3", "short_id3", 3)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", utils.BuildUrl("/stats", consts.GetPort()), nil)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Response code was %v; want 201", res.Code)
	}

	var result map[string]interface{}

	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&result)

	if result["totalClicks"].(float64) != 6 {
		t.Errorf("Total clicks was %v; want 6", result["totalClicks"])
	}

	if result["totalPastes"].(float64) != 3 {
		t.Errorf("Total pastes was %v; want 2", result["totalPastes"])
	}

	if result["avgClicks"].(float64) != 2 {
		t.Errorf("Average clicks was %v; want 2", result["avgClicks"])
	}
}

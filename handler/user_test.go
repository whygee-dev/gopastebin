package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"gopastebin/consts"
	"gopastebin/db"
	"gopastebin/models"
	"gopastebin/service"
	"gopastebin/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/argon2"
)


func TestLoginHappyPath(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	email := "admin@gmail.com"
	rawPassword := "password"
	password := service.HashPassword(rawPassword)

	_, err := db.Exec("INSERT INTO user (email, password) VALUES (?, ?)", email, password)

	if err != nil {
		t.Fatal(err)
	}

	body := map[string]string{
		"email": email,
		"password": rawPassword,
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(body)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/login", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Response code was %v; want 200", res.Code)
	}
}

func TestLoginPasswordIncorrect(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	email := "admin@gmail.com"
	_, err := db.Exec("INSERT INTO user (email, password) VALUES (?, ?)", email, "password")

	if err != nil {
		t.Fatal(err)
	}

	body := map[string]string{
		"email": email,
		"password": "any",
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(body)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/login", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Errorf("Response code was %v; want 401", res.Code)
	}
}

func TestLoginNotFound(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	body := map[string]string{
		"email": "any",
		"password": "any",
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(body)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/login", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Errorf("Response code was %v; want 401", res.Code)
	}
}

func TestLoginInvalidBody(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	buffer := new(bytes.Buffer)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/login", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("Response code was %v; want 500", res.Code)
	}
}

func TestLoginInternalError(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	_, err := db.Exec(`DROP TABLE user`)

	if err != nil {
		t.Fatal(err)
	}

	body := map[string]string{
		"email": "any",
		"password": "any",
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(body)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/login", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("Response code was %v; want 500", res.Code)
	}
}

func TestRegisterHappyPath(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	email := "admin@gmail.com"
	rawPassword := "password"

	body := map[string]string{
		"email": email,
		"password": rawPassword,
	}
	buffer := new(bytes.Buffer)

	json.NewEncoder(buffer).Encode(body)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/signup", consts.GetPort()), buffer)

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

		t.Errorf("Response body was %v; want id", result)
	}
	
	row := db.QueryRow("SELECT id, email, password FROM user WHERE email = ?", email)

	if row.Err() == sql.ErrNoRows {
		t.Errorf("User not found")
	}

	var user models.User

	row.Scan(&user.ID, &user.Email, &user.Password)

	if user.Email != email {
		t.Errorf("Email was %v; want %v", user.Email, email)
	}

	salt, time, memory, threads, keyLen := consts.GetArgonOptions()
	attemptedPassword := argon2.Key([]byte(rawPassword), salt, time, memory, threads, keyLen)

	if string(attemptedPassword) != user.Password {
		t.Errorf("Password was %v; want %v", user.Password, string(attemptedPassword))
	}

}

func TestRegisterInvalidBody(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	buffer := new(bytes.Buffer)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/signup", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Response code was %v; want 400", res.Code)
	}
}

func TestRegisterInvalidEmail(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	email := "admin"
	rawPassword := "password"

	body := map[string]string{
		"email": email,
		"password": rawPassword,
	}
	buffer := new(bytes.Buffer)

	json.NewEncoder(buffer).Encode(body)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/signup", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Response code was %v; want 400", res.Code)
	}
}

func TestRegisterEmptyEmail(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	email := ""
	rawPassword := "password"

	body := map[string]string{
		"email": email,
		"password": rawPassword,
	}
	buffer := new(bytes.Buffer)

	json.NewEncoder(buffer).Encode(body)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/signup", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Response code was %v; want 400", res.Code)
	}
}

func TestRegisterPasswordLessThan8Chars(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	email := "admin@gmail.com"
	rawPassword := "pass"

	body := map[string]string{
		"email": email,
		"password": rawPassword,
	}
	buffer := new(bytes.Buffer)

	json.NewEncoder(buffer).Encode(body)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/signup", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Response code was %v; want 400", res.Code)
	}
}

func TestRegisterUserExists(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	email := "admin@gmail.com"
	rawPassword := "password"

	_, err := db.Exec("INSERT INTO user (email, password) VALUES (?, ?)", email, rawPassword)

	if err != nil {
		t.Fatal(err)
	}

	body := map[string]string{
		"email": email,
		"password": rawPassword,
	}
	buffer := new(bytes.Buffer)

	json.NewEncoder(buffer).Encode(body)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/signup", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusConflict {
		t.Errorf("Response code was %v; want 409", res.Code)
	}
}


func TestRegisterInteralError(t *testing.T) {
	db := db.CreateDb()
	router := mux.NewRouter()
	SetupUserRoutes(db, router)
	defer utils.TestTearDown(db)

	_, err := db.Exec(`DROP TABLE user`)

	if err != nil {
		t.Fatal(err)
	}

	body := map[string]string{
		"email": "admin@gmail.com",
		"password": "password",
	}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(body)

	req, err := http.NewRequest("POST", utils.BuildUrl("/user/signup", consts.GetPort()), buffer)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusConflict {
		t.Errorf("Response code was %v; want 409", res.Code)
	}
}


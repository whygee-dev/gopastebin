package handler

import (
	"database/sql"
	"encoding/json"
	"gopastebin/consts"
	"gopastebin/models"
	"log"
	"net/http"
	"regexp"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/argon2"
)

func SetupUserRoutes(db *sql.DB, router *mux.Router) {
	router.HandleFunc("/user/signup", Signup(db)).Methods("POST")
	router.HandleFunc("/user/login", Login(db)).Methods("POST")
}


func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Login")

		var body = models.LoginUser{}

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&body)

		if err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		row := db.QueryRow("SELECT id, email, password FROM user WHERE email = ?", body.Email)

		var user = models.User{}

		row.Scan(&user.ID, &user.Email, &user.Password)

		salt, time, memory, threads, keyLen := consts.GetArgonOptions()
		attemptedPassword := argon2.Key([]byte(body.Password), salt, time, memory, threads, keyLen)

		if string(attemptedPassword) != user.Password {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{}).SignedString(consts.GetSecret())

		if err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(&token)
	}
}

func Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Signup")

		var body = models.CreateUser{}

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&body)

		if body.Email == "" || body.Password == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		if len(body.Password) < 8 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(body.Email) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		if err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		row := db.QueryRow("SELECT id FROM user WHERE email = ?", body.Email)

		if row.Scan() != sql.ErrNoRows {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)

			return
		}

		salt, time, memory, threads, keyLen := consts.GetArgonOptions()
		hashedPassword := argon2.Key([]byte(body.Password), salt, time, memory, threads, keyLen)

		_, err = db.Exec("INSERT INTO user (email, password) VALUES (?, ?)", body.Email, hashedPassword)

		if err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

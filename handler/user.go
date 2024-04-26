package handler

import (
	"database/sql"
	"encoding/json"
	"gopastebin/models"
	"gopastebin/service"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
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

		user, cust_err, err := service.GetUserByEmail(db, body.Email)

		if cust_err != nil {	
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		if err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		if !service.VerifyPassword(body.Password, user.Password) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		token, err := service.CreateToken()

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

		cust_err, err := service.CreateUser(db, body)

		if cust_err != nil {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)

			return
		}

		if err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
			

		w.WriteHeader(http.StatusCreated)
	}
}

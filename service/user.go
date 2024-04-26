package service

import (
	"database/sql"
	"gopastebin/custom_errors"
	"gopastebin/models"
)

func CreateUser(db *sql.DB, data models.CreateUser) (*custom_errors.UserExists, error) {
	row := db.QueryRow("SELECT id FROM user WHERE email = ?", data.Email)

	if row.Scan() != sql.ErrNoRows {
		return &custom_errors.UserExists{}, nil
	}

	hashedPassword := HashPassword(data.Password)

	_, err := db.Exec("INSERT INTO user (email, password) VALUES (?, ?)", data.Email, hashedPassword)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, *custom_errors.UserNotFound, error) {
	row := db.QueryRow("SELECT id, email, password FROM user WHERE email = ?", email)

	var user models.User

	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return nil, &custom_errors.UserNotFound{}, nil
	}

	if err != nil {
		return nil, nil, err
	}

	return &user, nil, nil
}
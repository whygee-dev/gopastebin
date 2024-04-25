package models

type User struct {
	ID    			int        `json:"-"`
	Email 			string
	Password 		string     `json:"-"`
}


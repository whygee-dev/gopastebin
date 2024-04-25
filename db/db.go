package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDb () *sql.DB {
	db, err := sql.Open("sqlite3", "./gopastebin.db")

	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	    create table if not exists paste
        (
            id integer not null primary key, 
            content text not null, 
            short_id text not null,
            click_count integer not null, 
            created_at datetime default current_timestamp
        );

        create unique index if not exists paste_short_id_uindex
        on paste (short_id);

		create table if not exists user
		(
			id integer not null primary key,
			email text not null,
			password text not null,
			created_at datetime default current_timestamp
		);

		create unique index if not exists user_email_uindex
		on user (email);
	`
	_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)

        panic(err)
	}

	return db
}
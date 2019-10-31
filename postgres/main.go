package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "test"
)

type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	CreatedOn time.Time
	LastLogin time.Time
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// query and scan
	var userid, username, email string

	sql := `SELECT user_id, username, email FROM account WHERE user_id=$1;`
	row := db.QueryRow(sql, 3)
	err = row.Scan(&userid, &username, &email)

	//if err == sql.ErrNoRows {
	//fmt.Println("No rows were returned!")
	//}

	if err != nil {
		panic(err)
	}

	fmt.Println(userid, username, email)

	// query and scan into struct

	var user User

	sql = `SELECT * FROM account WHERE user_id=$1;`
	row = db.QueryRow(sql, 3)
	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedOn, &user.LastLogin)

	//if err == sql.ErrNoRows {
	//fmt.Println("No rows were returned!")
	//}

	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", user)
	fmt.Printf("%s\n", user.CreatedOn)
	fmt.Printf("%s\n", user.LastLogin)
}

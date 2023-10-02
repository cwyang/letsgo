package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cwyang/letsgo/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	// we just use init() func only
)

// application-wide dependency
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	notes    *mysql.NotesModel
}

func main() {
	password := os.Getenv("MYSQL_USERPASS")
	var dsn string
	if password != "" {
		dsn = fmt.Sprintf("user:%s@/notes?parseTime=true", password)
	} else {
		dsn = *flag.String("dsn", "user:pass@/notes?parseTime=true",
			"MySQL data source name")
	}
	addr := flag.String("addr", ":4000", "HTTP listen port")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		notes:    &mysql.NotesModel{DB: db},
	}

	// err := http.ListenAndServe(*addr, mux)
	// redirect http.Server log to standard error
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("server started on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// returns a sql.DB connection pool for a givent DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

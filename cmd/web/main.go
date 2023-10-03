package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	
	"github.com/cwyang/letsgo/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	// we just use init() func only
	"github.com/golangcollege/sessions"
)

// application-wide dependency
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	notes    *mysql.NotesModel
	templateCache map[string]*template.Template
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
	secret := flag.String("secret", "12345678901234567890123456789012", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session:  session,
		notes:    &mysql.NotesModel{DB: db},
		templateCache: templateCache,
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

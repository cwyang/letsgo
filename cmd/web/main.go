package main

import (
	"crypto/tls"
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
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	templateCache map[string]*template.Template
	notes         *mysql.NoteModel
	users         *mysql.UserModel
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
	session.Secure = true

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		templateCache: templateCache,
		notes:         &mysql.NoteModel{DB: db},
		users:         &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}
	// err := http.ListenAndServe(*addr, mux)
	// redirect http.Server log to standard error
	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// timeouts
		IdleTimeout:  time.Minute, // go default keepalive 180s
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("server started on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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

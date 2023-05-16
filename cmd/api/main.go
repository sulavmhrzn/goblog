package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	dsn  string
}
type application struct {
	infolog  *log.Logger
	errorlog *log.Logger
	cfg      config
}

func main() {
	app := application{
		infolog:  log.New(os.Stdout, "INFO\t", log.Ltime|log.Lshortfile),
		errorlog: log.New(os.Stdout, "ERROR\t", log.Ltime|log.Lshortfile),
	}
	flag.IntVar(&app.cfg.port, "port", 4000, "Port number to serve")
	flag.StringVar(&app.cfg.dsn, "dsn", "postgres://goblog:goblogpw@localhost/goblog", "Database DSN")
	flag.Parse()
	// TODO:
	_, err := OpenDB(&app.cfg)
	if err != nil {
		app.errorlog.Fatal(err.Error())
	}

	app.infolog.Println("server running on port: ", app.cfg.port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", app.cfg.port), app.router())
	if err != nil {
		app.errorlog.Fatal(err)
	}
}

func OpenDB(cfg *config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

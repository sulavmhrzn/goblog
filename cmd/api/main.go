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
	"github.com/sulavmhrzn/goblog/internal/data"
)

type config struct {
	port int
	dsn  string
}
type application struct {
	infolog  *log.Logger
	errorlog *log.Logger
	config   config
	models   data.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Port number to serve")
	flag.StringVar(&cfg.dsn, "dsn", "", "Database DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := OpenDB(&cfg)
	if err != nil {
		logger.Fatal(err.Error())
	}
	app := application{
		infolog:  log.New(os.Stdout, "INFO\t", log.Ltime|log.Lshortfile),
		errorlog: log.New(os.Stdout, "ERROR\t", log.Ltime|log.Lshortfile),
		config:   cfg,
		models:   data.NewModels(db),
	}

	app.infolog.Println("Database connection successfull")
	app.infolog.Println("server running on port: ", cfg.port)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.router(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = srv.ListenAndServeTLS("./tls/localhost.pem", "./tls/localhost-key.pem")
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

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sulavmhrzn/goblog/internal/data"
	"github.com/sulavmhrzn/goblog/internal/mailer"
)

type config struct {
	port int
	dsn  string
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}
type application struct {
	infolog  *log.Logger
	errorlog *log.Logger
	config   config
	models   data.Models
	mailer   mailer.Mailer
}

func main() {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}

	flag.IntVar(&cfg.port, "port", 4000, "Port number to serve")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("DB_DSN"), "Database DSN")
	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host to connect to")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 0, "SMTP port")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "SMTP sender")
	flag.Parse()

	if cfg.smtp.port == 0 {
		port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
		if err != nil {
			log.Fatal(err)
		}
		cfg.smtp.port = port
	}

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
		mailer:   mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
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

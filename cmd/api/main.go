package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type config struct {
	port int
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
	flag.Parse()

	router := httprouter.New()

	app.infolog.Println("server running on port: ", app.cfg.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.cfg.port), router)
	if err != nil {
		app.errorlog.Fatal(err)
	}
}

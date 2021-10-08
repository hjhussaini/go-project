package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hjhussaini/go-project/charge/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const version = "1.0.0"

type config struct {
	port	int
	db	struct {
		uri	string
	}
	charge	struct {
		max		int
		credit		int
		walletAPI	string
	}
}

type application struct {
	version		string
	config		config
	database	*models.Model
	infoLog		*log.Logger
	errorLog	*log.Logger
}

func (app *application) serve() error {
	server := &http.Server{
		Addr:			fmt.Sprintf(":%d", app.config.port),
		Handler:		app.routes(),
		IdleTimeout:		30 * time.Second,
		ReadTimeout:		10 * time.Second,
		ReadHeaderTimeout:	5 * time.Second,
		WriteTimeout:		5 * time.Second,
	}

	app.infoLog.Printf("Starting Charge API server on port %d\n", app.config.port)

	return server.ListenAndServe()
}

func main() {
	var cfg config

	cfg.port, _ = strconv.Atoi(os.Getenv("PORT"))
	cfg.db.uri = os.Getenv("DATABASE_URI")
	cfg.charge.walletAPI = os.Getenv("CHARGE_WALLET_API")
	cfg.charge.max, _ = strconv.Atoi(os.Getenv("MAX_CHARGE"))
	cfg.charge.credit, _ = strconv.Atoi(os.Getenv("CHARGE_CREDIT"))

	infoLog := log.New(os.Stdout, "INFO ", log.Ldate | log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR ", log.Ldate | log.Ltime | log.Lshortfile)

	database, err := gorm.Open("mysql", cfg.db.uri)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer database.Close()

	app := &application{
		version:	version,
		config:		cfg,
		database:	models.New(database),
		infoLog:	infoLog,
		errorLog:	errorLog,
	}

	if err := app.serve(); err != nil {
		errorLog.Fatal(err)
	}
}

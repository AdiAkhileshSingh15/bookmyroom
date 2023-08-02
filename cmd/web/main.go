package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AdiAkhileshSingh15/bookmyroom/internal/config"
	"github.com/AdiAkhileshSingh15/bookmyroom/internal/driver"
	"github.com/AdiAkhileshSingh15/bookmyroom/internal/handlers"
	"github.com/AdiAkhileshSingh15/bookmyroom/internal/helpers"
	"github.com/AdiAkhileshSingh15/bookmyroom/internal/models"
	"github.com/AdiAkhileshSingh15/bookmyroom/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
)

const portNumber = ":9090"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	fmt.Printf("Starting application on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	_ = srv.ListenAndServe()
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	inProduction := os.Getenv("PRODUCTION")
	useCache := os.Getenv("CACHE")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbPort := os.Getenv("DB_PORT")
	dbSSL := os.Getenv("SSL_MODE")

	if inProduction == "true" {
		app.InProduction = true
	} else {
		app.InProduction = false
	}

	if useCache == "true" {
		app.UseCache = true
	} else {
		app.UseCache = false
	}

	if dbName == "" || dbUser == "" || dbPass == "" {
		fmt.Println("Missing required environment variables")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database

	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPass, dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
		return nil, err
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}

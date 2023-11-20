package main

import (
	"expvar"
	"flag"
	"fmt"
	"github.com/jessicatarra/greenlight/internal/database"
	"github.com/jessicatarra/greenlight/internal/jsonlog"
	"github.com/jessicatarra/greenlight/internal/mailer"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	version string
	port    string
	env     string
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
	jwt struct {
		secret string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models database.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

// @title Greenlight API Docs
// @version 1.0.0
// @contact.name Jessica Tarra
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	var cfg config

	intPort, _ := strconv.Atoi(port)
	flag.IntVar(&cfg.port, "port", intPort, "API server port")
	flag.StringVar(&cfg.env, "env", env, "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv(
		"DATABASE_URL"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", smtpPort, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.StringVar(&cfg.jwt.secret, "jwt-secret", os.Getenv("JWT_SECRET"), "JWT secret")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := database.New(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime, true)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: database.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
		return
	}
}

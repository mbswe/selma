package selma

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// LoggingConfig holds the logging configuration
type LoggingConfig struct {
	Middleware string `json:"middleware"`
	System     string `json:"system"`
	Debug      string `json:"debug"`
	Error      string `json:"error"`
	Info       string `json:"info"`
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// Config holds the server configuration
type Config struct {
	Mode       string         `json:"mode"`
	ServerPort int            `json:"server_port"`
	Logging    LoggingConfig  `json:"logging"`
	Database   DatabaseConfig `json:"database"`
}

// App houses the router and configuration
type App struct {
	Router           *mux.Router
	Config           *Config
	MiddlewareLogger *log.Logger
	DebugLogger      *log.Logger
	ViewRenderer     *ViewRenderer
	DB               *gorm.DB
}

// NewApp initializes the App with the configuration and router
func NewApp(configPath string) *App {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Failed to close config file: %v", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	router := mux.NewRouter()

	app := &App{
		Router: router,
		Config: config,
	}

	app.setupLogging()
	app.initDatabase()

	viewRenderer := NewViewRenderer(config, app.DebugLogger)
	if err := viewRenderer.LoadTemplates("views"); err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}
	app.ViewRenderer = viewRenderer

	return app
}

// setupLogging sets up logging based on the configuration
func (app *App) setupLogging() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Set up middleware logger
	middlewareLogFile, err := os.OpenFile(filepath.Join(logDir, app.Config.Logging.Middleware), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open middleware log file: %v", err)
	}
	app.MiddlewareLogger = log.New(middlewareLogFile, "", log.LstdFlags)

	debugLogFile, err := os.OpenFile(filepath.Join(logDir, app.Config.Logging.Debug), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open debug log file: %v", err)
	}
	app.DebugLogger = log.New(debugLogFile, "", log.LstdFlags)

	// Create a multi-writer to write the system log to both the log file and the console
	systemLogFile, err := os.OpenFile(filepath.Join(logDir, app.Config.Logging.System), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open system log file: %v", err)
	}
	multiWriter := io.MultiWriter(systemLogFile, os.Stdout)
	log.SetOutput(multiWriter)
}

// initDatabase initializes the database connection
func (app *App) initDatabase() {
	dbConfig := app.Config.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	app.DB = db
}

// RunMigrations runs database migrations using GORM
func (app *App) RunMigrations(models ...interface{}) {
	if err := app.DB.AutoMigrate(models...); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

// StartServer starts the HTTP server using the configuration
func (app *App) StartServer() {
	addr := ":" + strconv.Itoa(app.Config.ServerPort)

	log.Printf("Server listening on %s", addr)

	http.ListenAndServe(addr, app.Router)
}

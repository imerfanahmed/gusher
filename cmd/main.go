package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/imerfanahmed/gusher/internal/cache"
	"github.com/imerfanahmed/gusher/internal/config"
	"github.com/imerfanahmed/gusher/internal/database"
	"github.com/imerfanahmed/gusher/internal/types"
	"github.com/imerfanahmed/gusher/internal/webhook"
	"github.com/imerfanahmed/gusher/internal/websocket"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Define the --migrate flag
	migrate := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	// Set log flags for better debugging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Initialize database
	db, err := database.New(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	defer db.Close()

	// Run migrations only if --migrate flag is set
	if *migrate {
		if err := database.Migrate(db); err != nil {
			log.Fatal("Migration failed: ", err)
		}
	}

	// Initialize Redis
	redisClient, err := cache.NewRedisClient(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatal("Redis connection failed: ", err)
	}

	// Load configurations
	config.LoadAppConfigs(db, redisClient)
	webhook.LoadWebhookConfigs(db)

	// Set up HTTP router
	router := mux.NewRouter()
	router.HandleFunc("/app/{key}", websocket.Handler(db, redisClient)).Methods("GET")
	// Add other routes as needed

	// Start the server
	log.Printf("Server starting on %s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	if err := http.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), router); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
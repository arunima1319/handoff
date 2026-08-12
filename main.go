package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/arunima1319/handoff/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	dbQueries *database.Queries
	db        *sql.DB
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error in opening the database: %s", err)
		return
	}

	apiCfg := apiConfig{}
	apiCfg.db = db
	apiCfg.dbQueries = database.New(db)

	const port = "8080"
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/domains", apiCfg.handlerCreateDomain)
	mux.HandleFunc("POST /api/domains/{domainID}/users", apiCfg.handlerAddUserToDomain)
	mux.HandleFunc("POST /api/tasks", apiCfg.handlerCreateTask)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port %s...", port)
	log.Fatal(srv.ListenAndServe())
}

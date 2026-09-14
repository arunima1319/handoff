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
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error in opening the database: %s", err)
		return
	}

	apiCfg := apiConfig{}
	apiCfg.db = db
	apiCfg.dbQueries = database.New(db)

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir("./app")))
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/domains", apiCfg.handlerCreateDomain)
	mux.HandleFunc("POST /api/domains/{domainID}/users", apiCfg.handlerAddUserToDomain)
	//have to change this url path to "POST /api/domains/{domainID}/tasks"
	mux.HandleFunc("POST /api/tasks", apiCfg.handlerCreateTask)
	mux.HandleFunc("POST /api/tasks/{taskID}/dependencies", apiCfg.handlerCreateTaskDependency)
	mux.HandleFunc("GET /api/domains/{domainID}/tasks", apiCfg.handlerGetTasksOfDomain)
	mux.HandleFunc("GET /api/domains/{domainID}/users", apiCfg.handlerGetUsersOfDomain)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port %s...", port)
	log.Fatal(srv.ListenAndServe())
}

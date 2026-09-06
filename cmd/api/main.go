package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"movies-api/internal/database"
	"movies-api/internal/handler"
	"movies-api/internal/repository"
	"movies-api/internal/service"
)

func main() {

	fmt.Println("movies-api")

	err := os.MkdirAll("./data", 0755)
	if err != nil {
		log.Fatalf(">> FATAL ERROR: Failed to create directory \n%v", err)
	}

	dbPath := "./data/movies.db"
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf(">> ERROR: %v", err)
	}

	defer db.Close()

	migrationPath := "./migrations/0000001_init_schema.up.sql"
	err = database.RunMigration(db, migrationPath)
	if err != nil {
		log.Fatalf(">> FATAL ERROR: Could not run migration \n%v", err)
	}

	// Repositories
	genreRepo := repository.NewGenreRepository(db)
	actorRepo := repository.NewActorRepository(db)
	movieRepo := repository.NewMovieRepository(db)

	// Services
	genreService := service.NewGenreService(genreRepo)
	actorService := service.NewActorService(actorRepo)
	movieService := service.NewMovieService(movieRepo)

	// Handlers
	genreHandler := handler.NewGenreHandler(genreService)
	actorHandler := handler.NewActorHandler(actorService)
	movieHandler := handler.NewMovieHandler(movieService)

	port := ":8080"
	mux := http.NewServeMux()

	// Genres
	mux.HandleFunc("GET /health", HealthCheck)
	mux.HandleFunc("POST /api/genres", genreHandler.Create)
	mux.HandleFunc("GET /api/genres", genreHandler.GetAll)
	mux.HandleFunc("GET /api/genres/{id}", genreHandler.GetByID)
	mux.HandleFunc("PATCH /api/genres/{id}", genreHandler.Update)
	mux.HandleFunc("DELETE /api/genres/{id}", genreHandler.Delete)

	// Actors
	mux.HandleFunc("POST /api/actors", actorHandler.Create)
	mux.HandleFunc("GET /api/actors", actorHandler.GetAll)
	mux.HandleFunc("GET /api/actors/{id}", actorHandler.GetByID)
	mux.HandleFunc("PATCH /api/actors/{id}", actorHandler.Update)
	mux.HandleFunc("DELETE /api/actors/{id}", actorHandler.Delete)

	// Movies
	mux.HandleFunc("POST /api/movies", movieHandler.Create)
	mux.HandleFunc("GET /api/movies", movieHandler.GetAll)
	mux.HandleFunc("GET /api/movies/{id}", movieHandler.GetByID)

	mux.HandleFunc("GET /api/movies/search", movieHandler.Search)
	mux.HandleFunc("GET /api/movies/{id}/actors", movieHandler.GetActors)

	mux.HandleFunc("PATCH /api/movies/{id}", movieHandler.Update)
	mux.HandleFunc("DELETE /api/movies/{id}", movieHandler.Delete)

	fmt.Printf(">> Starting Server ...\n")
	fmt.Printf(">> URL: http://localhost%s\n", port)

	handlerWithMiddleware := PanicRecoveryMiddleware(mux)

	err = http.ListenAndServe(port, handlerWithMiddleware )
	if err != nil {
		log.Fatalf(">> FATAL: Server Failed to start! \n%v", err)
	}

}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(">> OK: API Health Check!"))
	fmt.Printf(">> OK: API Health Check!\n")
}


func PanicRecoveryMiddleware(next http.Handler ) http.Handler {
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request ) {
		
		defer func() {

			if err := recover(); err != nil {
				log.Printf(">> PANIC RECOVERED: %v\n", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError )
				w.Write([]byte(`{"error": "Internal server error"}`))
			}
		}()
		
		next.ServeHTTP(w, r)

	})

}
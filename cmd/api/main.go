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

	// fmt.Printf(">> TEST: Genre Repository Call...\n")

	// newGenre := &models.Genre{
	// 	Name: "Thriller",
	// }

	// fmt.Printf("\n--- Create Method ---\n")
	// fmt.Printf("Initial Value > newGenre: %+v\n", newGenre)

	// err = genreRepo.Create(newGenre)
	// if err != nil {
	// 	log.Printf("Error creating genre: %v\n", err)
	// } else {
	// 	fmt.Printf("Update Value > newGenre:  %+v\n", newGenre)
	// }
	// fmt.Printf("---------------------\n")

	// fmt.Printf("\n--- GetAll Method ---\n")
	// allGenres, err := genreRepo.GetAll()
	// if err != nil {
	// 	log.Printf("Error getting genre: %v\n", err)
	// } else {
	// 	fmt.Printf("Found %d genres: \n", len(allGenres))
	// 	for _, g := range allGenres {
	// 		fmt.Printf(" - ID: %d, Name: %s\n", g.ID, g.Name)
	// 	}
	// }
	// fmt.Printf("---------------------\n")

	// fmt.Printf("\n--- Update Method ---\n")
	// genreToUpdate := &models.Genre{
	// 	ID:   3,
	// 	Name: "Horror",
	// }
	// err = genreRepo.Update(genreToUpdate)
	// if err != nil {
	// 	log.Printf("Error updating genre: %v\n", err)
	// } else {
	// 	fmt.Println("Successfully updated genre ID'!")
	// }

	// fmt.Printf("---------------------\n")

	// fmt.Printf("\n--- Delete Method ---\n")
	// err = genreRepo.Delete(3)
	// if err != nil {
	// 	log.Printf("Error deleting genre: %v\n", err)
	// } else {
	// 	fmt.Println("Successfully deleted genre ID!")
	// }

	// fmt.Printf("---------------------\n")

	// fmt.Printf("\n--- Verification Method ---\n")
	// genre, err := genreRepo.GetByID(2)
	// if err != nil {
	// 	fmt.Printf("Verification: %v\n", err)
	// } else {
	// 	fmt.Printf(" - ID: %d, Name: %s\n", genre.ID, genre.Name)
	// }
	// fmt.Printf("---------------------\n")

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

	// IMPORTANT: Keep it above {id} else it will break
	mux.HandleFunc("GET /api/movies/search", movieHandler.Search)

	mux.HandleFunc("GET /api/movies/{id}", movieHandler.GetByID)
	mux.HandleFunc("PATCH /api/movies/{id}", movieHandler.Update)
	mux.HandleFunc("DELETE /api/movies/{id}", movieHandler.Delete)

	fmt.Printf(">> Starting Server ...\n")
	fmt.Printf(">> URL: http://localhost%s\n", port)

	err = http.ListenAndServe(port, mux)
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
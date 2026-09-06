package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"movies-api/internal/models"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {

	return &MovieRepository{
		db: db,
	}
}

func (r *MovieRepository) Create(movie *models.Movie) error {

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO movies (title, release_year, duration) VALUES (?,?,?)`

	result, err := tx.Exec(query, movie.Title, movie.ReleaseYear, movie.Duration)
	if err != nil {
		return fmt.Errorf("failed to enter movie: %w", err)
	}

	movieID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get movieID: %w", err)
	}

	movie.ID = movieID

	for _, genre := range movie.Genres {

		genreQuery := `INSERT INTO movie_genres (movie_id, genre_id) VALUES (?, ?)`
		_, err := tx.Exec(genreQuery, movie.ID, genre.ID)
		if err != nil {
			return fmt.Errorf("failed to link genre ID %d: %w", genre.ID, err)
		}

	}

	for _, actor := range movie.Actors {

		actorQuery := `INSERT INTO movie_actors (movie_id, actor_id) VALUES (?, ?)`
		_, err := tx.Exec(actorQuery, movie.ID, actor.ID)
		if err != nil {
			return fmt.Errorf("failed to link actor ID %d: %w", actor.ID, err)
		}

	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil

}

func (r *MovieRepository) GetAll() ([]models.Movie, error) {

	movieQuery := `SELECT id, title, release_year, duration FROM movies`
	movieRows, err := r.db.Query(movieQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query movies: %w", err)
	}

	defer movieRows.Close()

	movieMap := make(map[int64]*models.Movie)

	var movieIDs []int64

	for movieRows.Next() {

		var m models.Movie

		err := movieRows.Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}

		m.Genres = []models.Genre{}
		m.Actors = []models.Actor{}

		movieMap[m.ID] = &m
		movieIDs = append(movieIDs, m.ID)

	}

	if len(movieMap) == 0 {
		return []models.Movie{}, nil
	}

	genreQuery := `
		SELECT movie_genres.movie_id, genres.id, genres.name
		FROM genres
		INNER JOIN movie_genres ON genres.id = movie_genres.genre_id
	`
	genreRows, err := r.db.Query(genreQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query movie genres: %w", err)
	}
	defer genreRows.Close()

	for genreRows.Next() {

		var movieID int64
		var g models.Genre

		err := genreRows.Scan(&movieID, &g.ID, &g.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie genre: %w", err)
		}

		movie, exists := movieMap[movieID]
		if exists == true {
			movie.Genres = append(movie.Genres, g)
		}
	}

	actorQuery := `
		SELECT movie_actors.movie_id, actors.id, actors.name, actors.birth_date
		FROM actors
		INNER JOIN movie_actors ON actors.id = movie_actors.actor_id
	`
	actorRows, err := r.db.Query(actorQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query movie actors: %w", err)
	}
	defer actorRows.Close()

	for actorRows.Next() {

		var movieID int64
		var a models.Actor

		err := actorRows.Scan(&movieID, &a.ID, &a.Name, &a.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie actor: %w", err)
		}

		movie, exists := movieMap[movieID]
		if exists == true {
			movie.Actors = append(movie.Actors, a)
		}
	}

	var finalMovies []models.Movie

	for _, id := range movieIDs {
		finalMovies = append(finalMovies, *movieMap[id])
	}

	return finalMovies, nil
}

func (r *MovieRepository) GetByID(id int64) (models.Movie, error) {

	var m models.Movie
	m.Genres = []models.Genre{}
	m.Actors = []models.Actor{}

	movieQuery := `SELECT id, title, release_year, duration FROM movies WHERE id = ?`

	err := r.db.QueryRow(movieQuery, id).Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m, fmt.Errorf("movie with ID %d not found", id)
		}
		return m, fmt.Errorf("failed to query movie: %w", err)
	}

	genreQuery := `
		SELECT genres.id, genres.name
		FROM genres
		INNER JOIN movie_genres movie_genres ON genres.id = movie_genres.genre_id
		WHERE movie_genres.movie_id = ?
	`

	genreRows, err := r.db.Query(genreQuery, id)
	if err != nil {
		return m, fmt.Errorf("failed to query movie genres: %w", err)
	}
	defer genreRows.Close()

	for genreRows.Next() {
		var g models.Genre

		err := genreRows.Scan(&g.ID, &g.Name)
		if err != nil {
			return m, fmt.Errorf("failed to scan genre: %w", err)
		}

		m.Genres = append(m.Genres, g)

	}

	actorQuery := `
		SELECT actors.id, actors.name, actors.birth_date
		FROM actors
		INNER JOIN movie_actors movie_actors ON actors.id = movie_actors.actor_id
		WHERE movie_actors.movie_id = ?
	`

	actorRows, err := r.db.Query(actorQuery, id)
	if err != nil {
		return m, fmt.Errorf("failed to query movie actors: %w", err)
	}
	defer actorRows.Close()

	for actorRows.Next() {
		var a models.Actor

		err := actorRows.Scan(&a.ID, &a.Name, &a.BirthDate)
		if err != nil {
			return m, fmt.Errorf("failed to scan actor: %w", err)
		}

		m.Actors = append(m.Actors, a)

	}

	return m, nil
}

func (r *MovieRepository) Update(movie *models.Movie) error {

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `UPDATE movies SET title = ?, release_year = ?, duration = ? WHERE id = ?`

	_, err = tx.Exec(query, movie.Title, movie.ReleaseYear, movie.Duration, movie.ID)
	if err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM movie_genres WHERE movie_id = ?`, movie.ID)
	if err != nil {
		return fmt.Errorf("failed to delete old genres: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM movie_actors WHERE movie_id = ?`, movie.ID)
	if err != nil {
		return fmt.Errorf("failed to delete old actors: %w", err)
	}

	for _, genre := range movie.Genres {

		_, err = tx.Exec(`INSERT INTO movie_genres (movie_id, genre_id) VALUES (?, ?)`, movie.ID, genre.ID)
		if err != nil {
			return fmt.Errorf("failed to insert new genre: %w", err)
		}

	}

	for _, actor := range movie.Actors {

		_, err = tx.Exec(`INSERT INTO movie_actors (movie_id, actor_id) VALUES (?, ?)`, movie.ID, actor.ID)
		if err != nil {
			return fmt.Errorf("failed to insert new actor: %w", err)
		}

	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *MovieRepository) Delete(id int64) error {

	query := `DELETE FROM movies WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	return nil
}

func (r *MovieRepository) ForceDelete(id int64) error {

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM movie_genres WHERE movie_id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie genres: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM movie_actors WHERE movie_id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie actors: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM movies WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *MovieRepository) SearchByTitle(title string) ([]models.Movie, error) {

	query := `SELECT id, title, release_year, duration FROM movies WHERE title LIKE ?`

	searchPattern := "%" + title + "%"

	rows, err := r.db.Query(query, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}

		m.Genres = []models.Genre{}
		m.Actors = []models.Actor{}
		movies = append(movies, m)
	}

	return movies, nil
}


func (r *MovieRepository) GetByActorID(actorID int64) ([]models.Movie, error) {
	
	query := `
		SELECT m.id, m.title, m.release_year, m.duration 
		FROM movies m
		INNER JOIN movie_actors ma ON m.id = ma.movie_id
		WHERE ma.actor_id = ?
	`

	rows, err := r.db.Query(query, actorID)
	
	if err != nil {
		return nil, fmt.Errorf("failed to query movies by actor ID: %w", err)
	}
	defer rows.Close()

	var movies []models.Movie

	for rows.Next() {
	
		var m models.Movie
	
		err := rows.Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}
	
		m.Genres = []models.Genre{}
		m.Actors = []models.Actor{}
	
		movies = append(movies, m)
	
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if movies == nil {
		return []models.Movie{}, nil
	}

	return movies, nil
}


func (r *MovieRepository) GetByGenreID(genreID int64) ([]models.Movie, error) {
	
	query := `
		SELECT m.id, m.title, m.release_year, m.duration 
		FROM movies m
		INNER JOIN movie_genres mg ON m.id = mg.movie_id
		WHERE mg.genre_id = ?
	`

	rows, err := r.db.Query(query, genreID)
	
	if err != nil {
		return nil, fmt.Errorf("failed to query movies by genre ID: %w", err)
	}
	
	defer rows.Close()

	var movies []models.Movie
	
	for rows.Next() {
	
		var m models.Movie
	
		err := rows.Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}
	
		m.Genres = []models.Genre{}
		m.Actors = []models.Actor{}
	
		movies = append(movies, m)
	
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if movies == nil {
		return []models.Movie{}, nil
	}

	return movies, nil
}


func (r *MovieRepository) GetByReleaseYear(year int) ([]models.Movie, error) {
	
	query := `SELECT id, title, release_year, duration FROM movies WHERE release_year = ?`

	rows, err := r.db.Query(query, year)
	
	if err != nil {
		return nil, fmt.Errorf("failed to query movies by release year: %w", err)
	}
	
	defer rows.Close()

	var movies []models.Movie
	
	for rows.Next() {
	
		var m models.Movie
	
		err := rows.Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}
	
		m.Genres = []models.Genre{}
		m.Actors = []models.Actor{}
	
		movies = append(movies, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if movies == nil {
		return []models.Movie{}, nil
	}

	return movies, nil
}
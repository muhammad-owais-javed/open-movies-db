package repository

import (
	"fmt"
	"errors"
	"database/sql"
	"movies-api/internal/models"
)

type GenreRepository struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{
		db: db,
	}
}

func (r *GenreRepository) Create(genre *models.Genre) error {

	query := `INSERT INTO genres (name) VALUES (?)`

	result, err := r.db.Exec(query, genre.Name)
	if err != nil {
		fmt.Errorf("failed to enter genre: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Errorf("failed to get new ID: %w", err)
	}

	genre.ID = id
	

	return nil
}

func (r *GenreRepository) GetAll() ([]models.Genre, error){

	query := `SELECT id, name FROM genres`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query genres: %w", err)
	}

	defer rows.Close()

	var genres []models.Genre

	for rows.Next(){

		var g models.Genre

		err := rows.Scan(&g.ID, &g.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan genres: %w", err)
		}

		genres = append(genres, g)

	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row interaction error: %w", err)
	}

	return genres, nil

}

func (r *GenreRepository) GetByID(id int64) (*models.Genre, error) {

	query := `SELECT id, name FROM genres WHERE id = ?`

	var g models.Genre 

	err := r.db.QueryRow(query, id).Scan(&g.ID, &g.Name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("genre with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to query genre: %w", err)
	}

	return &g, nil
}

func (r *GenreRepository) GetByName(name string) ([]models.Genre, error) {
	
	query := `SELECT id, name FROM genres WHERE name LIKE ?`
	
	searchPattern := "%" + name + "%"

	rows, err := r.db.Query(query, searchPattern)
	
	if err != nil {
		return nil, fmt.Errorf("failed to search genres by name: %w", err)
	}
	
	defer rows.Close()

	var genres []models.Genre

	for rows.Next() {
	
		var g models.Genre
		err := rows.Scan(&g.ID, &g.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan genre: %w", err)
		}
		genres = append(genres, g)
	
	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if genres == nil {
		return []models.Genre{}, nil
	}

	return genres, nil
}


func (r *GenreRepository) Update(genre *models.Genre) error {
	query := `UPDATE genres SET name = ? WHERE id = ?`
	
result, err := r.db.Exec(query, genre.Name, genre.ID)
	if err != nil {
		return fmt.Errorf("failed to update genre: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("genre with ID %d not found", genre.ID)
	}

	return nil
}

func (r *GenreRepository) Delete(id int64) error {
	
	query := `DELETE FROM genres WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete genre: %w", err)
	}

	return nil

}
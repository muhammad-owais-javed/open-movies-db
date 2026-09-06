package repository

import (
	"fmt"
	"errors"
	"database/sql"
	"movies-api/internal/models"
)

type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository{
	return &ActorRepository{
		db: db,
	}
}

func (r *ActorRepository) Create(actor *models.Actor) error{

	query := `INSERT INTO actors (name, birth_date) VALUES (?, ?)`

	result, err := r.db.Exec(query, actor.Name, actor.BirthDate)
	if err != nil {
		return fmt.Errorf("failed to enter actor: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get new ID: %w", err)
	}

	actor.ID = id

	return nil
}

func (r *ActorRepository) GetAll()([]models.Actor, error){

	query := `SELECT id, name, birth_date FROM actors`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query actors: %w", err)
	}

	defer rows.Close()

	var actors []models.Actor
	
	for rows.Next(){
		
		var a models.Actor

		err := rows.Scan(&a.ID, &a.Name, &a.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan actor: %w", err)
		}

		actors = append(actors, a)
	}
	
	err = rows.Err();
	if err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return actors, nil
}

func (r *ActorRepository) GetByID(id int64) (*models.Actor, error) {

	query := `SELECT id, name, birth_date FROM actors WHERE id = ?`

	var a models.Actor

	err := r.db.QueryRow(query,id).Scan(&a.ID, &a.Name, &a.BirthDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("actor with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to query actor: %w", err)
	}

	return &a, nil

}


func (r *ActorRepository) GetByName(name string) ([]models.Actor, error) {
	
	query := `SELECT id, name, birth_date FROM actors WHERE name LIKE ?`
	
	searchPattern := "%" + name + "%"

	rows, err := r.db.Query(query, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search actors by name: %w", err)
	}

	defer rows.Close()

	var actors []models.Actor

	for rows.Next() {

		var a models.Actor
		err := rows.Scan(&a.ID, &a.Name, &a.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan actor: %w", err)
		}
		actors = append(actors, a)

	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if actors == nil {
		return []models.Actor{}, nil
	}

	return actors, nil
}


func (r *ActorRepository) Update(actor *models.Actor) error{

	query := `UPDATE actors SET name = ?, birth_date = ? WHERE id = ?`
	
	result, err := r.db.Exec(query, actor.Name, actor.BirthDate, actor.ID)

	if err != nil {
		return fmt.Errorf("failed to update actor: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("actor with ID %d not found", actor.ID) 
	}

	return nil

}


func (r *ActorRepository) Delete(id int64) error{

	query := `DELETE FROM actors WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete actor: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("actor with ID %d not found", id) 
	}

	return nil
}

func (r *ActorRepository) ForceDelete(id int64) error{

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM movie_actors WHERE actor_id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete actor relationships: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM actors WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete actor: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil


}


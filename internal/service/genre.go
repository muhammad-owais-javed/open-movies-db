package service

import (
	"fmt"
	"strings"
	"errors"
	"movies-api/internal/repository"
	"movies-api/internal/models"
)

type GenreService struct {
	repo *repository.GenreRepository
}

func NewGenreService(repo *repository.GenreRepository) *GenreService {
	return &GenreService{
		repo: repo,
	}
}

func (s *GenreService) Create(genre *models.Genre) error {
	
	genre.Name = strings.TrimSpace(genre.Name)
	if genre.Name == "" {
		return errors.New("genre name cannot be empty")
	} 
	
	err := s.repo.Create(genre)
	if err != nil {
		return fmt.Errorf("service failed to create genre: %w", err)
	}

	return nil
}

func (s *GenreService) GetAll() ([]models.Genre, error) {
	return s.repo.GetAll()
}

func (s *GenreService) GetByID(id int64) (*models.Genre, error) {
	if id <= 0 {
		return nil, errors.New("invalid ID provided")
	}
	return s.repo.GetByID(id)
}

func (s *GenreService) GetByName(name string) ([]models.Genre, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("search name cannot be empty")
	}
	return s.repo.GetByName(name)
}


func (s *GenreService) Update(genre *models.Genre) error {
	
	genre.Name = strings.TrimSpace(genre.Name)
	
	if genre.Name == "" {
		return errors.New("genre name cannot be empty")
	}
	if genre.ID <= 0 {
		return errors.New("invalid ID provided")
	}

	return s.repo.Update(genre)
}


func (s *GenreService) Delete(id int64, force bool) error {
	
	if id <= 0 {
		return errors.New("invalid ID provided")
	}

	if force == true {
		fmt.Println("Force delete requested!")
	}

	return s.repo.Delete(id)
}
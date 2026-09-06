package service

import (
	"errors"
	"fmt"
	"strings"

	"movies-api/internal/models"
	"movies-api/internal/repository"
)

type MovieService struct {
	repo *repository.MovieRepository
}

func NewMovieService(repo *repository.MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) Create(movie *models.Movie) error {

	movie.Title = strings.TrimSpace(movie.Title)
	if movie.Title == "" {
		return errors.New("movie title cannot be empty")
	}

	if movie.ReleaseYear < 1888 {
		return errors.New("release year must be 1888 or later")
	}

	if movie.Duration <= 0 {
		return errors.New("duration must be greater than 0 minutes")
	}

	err := s.repo.Create(movie)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return errors.New("the provided genre ID or actor ID does not exist")
		}
		return fmt.Errorf("failed to create movie: %w", err)
	}

	return nil
}

func (s *MovieService) GetAll() ([]models.Movie, error) {

	return s.repo.GetAll()

}

func (s *MovieService) GetByID(id int64) (models.Movie, error) {

	if id <= 0 {
		return models.Movie{}, errors.New("invalid ID provided")
	}

	return s.repo.GetByID(id)
}

func (s *MovieService) Update(id int64, updateData *models.Movie) error {

	if id <= 0 {
		return errors.New("invalid ID provided")
	}

	existingMovie, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if strings.TrimSpace(updateData.Title) != "" {
		existingMovie.Title = strings.TrimSpace(updateData.Title)
	}

	if updateData.ReleaseYear > 0 {
		if updateData.ReleaseYear < 1888 {
			return errors.New("release year must be 1888 or later")
		}

		existingMovie.ReleaseYear = updateData.ReleaseYear
	}

	if updateData.Duration > 0 {
		existingMovie.Duration = updateData.Duration
	}

	// For relationships (Genres/Actors), if user provide new list, we replace the old one.
	if updateData.Genres != nil {
		existingMovie.Genres = updateData.Genres
	}
	if updateData.Actors != nil {
		existingMovie.Actors = updateData.Actors
	}

	return s.repo.Update(&existingMovie)
}

func (s *MovieService) Delete(id int64, force bool) error {

	if id <= 0 {
		return errors.New("invalid ID provided")
	}

	if force == true {
		return s.repo.ForceDelete(id)
	}

	return s.repo.Delete(id)
}

func (s *MovieService) SearchByTitle(title string) ([]models.Movie, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("search title cannot be empty")
	}

	return s.repo.SearchByTitle(title)
}

func (s *MovieService) GetByActorID(actorID int64) ([]models.Movie, error) {
	if actorID <= 0 {
		return nil, errors.New("invalid actor ID provided")
	}
	return s.repo.GetByActorID(actorID)
}

func (s *MovieService) GetByGenreID(genreID int64) ([]models.Movie, error) {
	if genreID <= 0 {
		return nil, errors.New("invalid genre ID provided")
	}
	return s.repo.GetByGenreID(genreID)
}

func (s *MovieService) GetByReleaseYear(year int) ([]models.Movie, error) {
	if year < 1888 {
		return nil, errors.New("release year must be 1888 or later")
	}
	return s.repo.GetByReleaseYear(year)
}
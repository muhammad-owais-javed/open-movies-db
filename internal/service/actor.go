package service

import (
	"fmt"
	"strings"
	"errors"
	"time"
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

type ActorService struct{
	repo *repository.ActorRepository
}

func NewActorService(repo *repository.ActorRepository) *ActorService{
	return &ActorService{
		repo: repo,
	}
}

func (s *ActorService) Create(actor *models.Actor) error{

	actor.Name = strings.TrimSpace(actor.Name)
	if actor.Name == "" {
		return errors.New("actor name cannot be empty")
	}

	_, err := time.Parse("2006-01-02", actor.BirthDate)
	if err != nil {
		return errors.New("birth date must be in YYYY-MM-DD format")
	}

	err = s.repo.Create(actor)
	if err != nil {
		return fmt.Errorf("failed to create actor: %w", err)
	}

	return nil

}

func (s *ActorService) GetAll() ([]models.Actor, error) {
	return s.repo.GetAll()
}

func (s *ActorService) GetByID(id int64) (*models.Actor, error) {
	
	if id <= 0 {
		return &models.Actor{}, errors.New("invalid ID provided")
	}

	return s.repo.GetByID(id)

}

func (s *ActorService) GetByName(name string) ([]models.Actor, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("search name cannot be empty")
	}
	return s.repo.GetByName(name)
}


func (s *ActorService) Update(id int64, updateData *models.Actor) error {

	if id <= 0 {
		return errors.New("invalid ID provided")
	}

	existingActor, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if strings.TrimSpace(updateData.Name) != "" {
		existingActor.Name = strings.TrimSpace(updateData.Name)
	}
	
	if updateData.BirthDate != "" {
		_, err := time.Parse("2006-01-02", updateData.BirthDate) 
		if err != nil {
			return errors.New("birth date must be in YYYY-MM-DD format")
		}
		existingActor.BirthDate = updateData.BirthDate
	}

	return s.repo.Update(existingActor)
}


func (s *ActorService) Delete(id int64, force bool) error {
	
	if id <= 0 {
		return errors.New("invalid ID provided")
	}

	if force == true {
		return s.repo.ForceDelete(id)
	}

	return s.repo.Delete(id)
}
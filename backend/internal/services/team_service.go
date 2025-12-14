package services

import (
	"github.com/zahidcakici/champions-league/internal/models"
	"github.com/zahidcakici/champions-league/internal/repository"
)

type TeamService interface {
	GetAllTeams() ([]models.Team, error)
	CreateTeam(name string, power int) error
	DeleteTeam(id uint) error
	SeedTeams() error
}

type teamService struct {
	teamRepo repository.TeamRepository
}

func NewTeamService(teamRepo repository.TeamRepository) TeamService {
	return &teamService{teamRepo: teamRepo}
}

func (s *teamService) GetAllTeams() ([]models.Team, error) {
	// Auto-seed if no teams exist
	if err := s.teamRepo.SeedDefault(); err != nil {
		return nil, err
	}
	return s.teamRepo.FindAll()
}

func (s *teamService) CreateTeam(name string, power int) error {
	team := &models.Team{
		Name:  name,
		Power: power,
	}
	return s.teamRepo.Create(team)
}

func (s *teamService) DeleteTeam(id uint) error {
	return s.teamRepo.Delete(id)
}

func (s *teamService) SeedTeams() error {
	return s.teamRepo.SeedDefault()
}

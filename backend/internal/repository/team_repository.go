package repository

import (
	"github.com/zahidcakici/champions-league/internal/models"
	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(team *models.Team) error
	FindAll() ([]models.Team, error)
	FindByID(id uint) (*models.Team, error)
	FindByName(name string) (*models.Team, error)
	Count() (int64, error)
	Delete(id uint) error
	DeleteAll() error
	SeedDefault() error
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

func (r *teamRepository) FindAll() ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Order("id").Find(&teams).Error
	return teams, err
}

func (r *teamRepository) FindByID(id uint) (*models.Team, error) {
	var team models.Team
	err := r.db.First(&team, id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) FindByName(name string) (*models.Team, error) {
	var team models.Team
	err := r.db.Where("name = ?", name).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Team{}).Count(&count).Error
	return count, err
}

func (r *teamRepository) Delete(id uint) error {
	return r.db.Delete(&models.Team{}, id).Error
}

func (r *teamRepository) DeleteAll() error {
	return r.db.Exec("DELETE FROM teams").Error
}

func (r *teamRepository) SeedDefault() error {
	count, err := r.Count()
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Already seeded
	}

	teams := models.DefaultTeams()
	for _, team := range teams {
		if err := r.Create(&team); err != nil {
			return err
		}
	}
	return nil
}

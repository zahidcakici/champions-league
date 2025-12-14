package repository

import (
	"github.com/zahidcakici/champions-league/internal/models"
	"gorm.io/gorm"
)

type LeagueStateRepository interface {
	Get() (*models.LeagueState, error)
	Create(state *models.LeagueState) error
	Update(state *models.LeagueState) error
	Reset() error
}

type leagueStateRepository struct {
	db *gorm.DB
}

func NewLeagueStateRepository(db *gorm.DB) LeagueStateRepository {
	return &leagueStateRepository{db: db}
}

func (r *leagueStateRepository) Get() (*models.LeagueState, error) {
	var state models.LeagueState
	err := r.db.First(&state).Error
	if err == gorm.ErrRecordNotFound {
		// Create default state
		state = models.LeagueState{
			CurrentWeek:     0,
			TotalWeeks:      6,
			FixturesCreated: false,
			Started:         false,
			Completed:       false,
		}
		if createErr := r.db.Create(&state).Error; createErr != nil {
			return nil, createErr
		}
		return &state, nil
	}
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *leagueStateRepository) Create(state *models.LeagueState) error {
	return r.db.Create(state).Error
}

func (r *leagueStateRepository) Update(state *models.LeagueState) error {
	return r.db.Save(state).Error
}

func (r *leagueStateRepository) Reset() error {
	return r.db.Exec("DELETE FROM league_states").Error
}

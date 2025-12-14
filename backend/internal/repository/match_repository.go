package repository

import (
	"github.com/zahidcakici/champions-league/internal/models"
	"gorm.io/gorm"
)

type MatchRepository interface {
	Create(match *models.Match) error
	CreateBatch(matches []models.Match) error
	FindAll() ([]models.Match, error)
	FindByID(id uint) (*models.Match, error)
	FindByWeek(week int) ([]models.Match, error)
	FindPlayedMatches() ([]models.Match, error)
	Update(match *models.Match) error
	DeleteAll() error
	GetMaxWeek() (int, error)
}

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) Create(match *models.Match) error {
	return r.db.Create(match).Error
}

func (r *matchRepository) CreateBatch(matches []models.Match) error {
	return r.db.Create(&matches).Error
}

func (r *matchRepository) FindAll() ([]models.Match, error) {
	var matches []models.Match
	err := r.db.Preload("HomeTeam").Preload("AwayTeam").Order("week, id").Find(&matches).Error
	return matches, err
}

func (r *matchRepository) FindByID(id uint) (*models.Match, error) {
	var match models.Match
	err := r.db.Preload("HomeTeam").Preload("AwayTeam").First(&match, id).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) FindByWeek(week int) ([]models.Match, error) {
	var matches []models.Match
	err := r.db.Preload("HomeTeam").Preload("AwayTeam").Where("week = ?", week).Order("id").Find(&matches).Error
	return matches, err
}

func (r *matchRepository) FindPlayedMatches() ([]models.Match, error) {
	var matches []models.Match
	err := r.db.Preload("HomeTeam").Preload("AwayTeam").Where("played = ?", true).Order("week, id").Find(&matches).Error
	return matches, err
}

func (r *matchRepository) Update(match *models.Match) error {
	return r.db.Save(match).Error
}

func (r *matchRepository) DeleteAll() error {
	return r.db.Exec("DELETE FROM matches").Error
}

func (r *matchRepository) GetMaxWeek() (int, error) {
	var maxWeek int
	err := r.db.Model(&models.Match{}).Select("COALESCE(MAX(week), 0)").Scan(&maxWeek).Error
	return maxWeek, err
}

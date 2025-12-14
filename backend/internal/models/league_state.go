package models

import (
	"time"
)

type LeagueState struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CurrentWeek     int       `json:"current_week" gorm:"default:0"`
	TotalWeeks      int       `json:"total_weeks" gorm:"default:6"`
	FixturesCreated bool      `json:"fixtures_created" gorm:"default:false"`
	Started         bool      `json:"started" gorm:"default:false"`
	Completed       bool      `json:"completed" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

package models

import (
	"time"
)

type Team struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex;not null"`
	Power     int       `gorm:"not null;default:50"` // Team strength 1-100
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// DefaultTeams returns the 4 seeded teams with their power ratings
func DefaultTeams() []Team {
	return []Team{
		{Name: "Chelsea", Power: 85},
		{Name: "Arsenal", Power: 80},
		{Name: "Manchester City", Power: 90},
		{Name: "Liverpool", Power: 82},
	}
}

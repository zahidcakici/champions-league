package models

import (
	"time"
)

type Match struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Week       int       `json:"week" gorm:"not null;index"`
	HomeTeamID uint      `json:"home_team_id" gorm:"not null"`
	AwayTeamID uint      `json:"away_team_id" gorm:"not null"`
	HomeScore  *int      `json:"home_score"` // nil if not played
	AwayScore  *int      `json:"away_score"` // nil if not played
	Played     bool      `json:"played" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	HomeTeam Team `json:"home_team" gorm:"foreignKey:HomeTeamID"`
	AwayTeam Team `json:"away_team" gorm:"foreignKey:AwayTeamID"`
}

type MatchResult struct {
	HomeTeamName string `json:"home_team_name"`
	AwayTeamName string `json:"away_team_name"`
	HomeScore    int    `json:"home_score"`
	AwayScore    int    `json:"away_score"`
}

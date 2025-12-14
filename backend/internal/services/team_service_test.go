package services

import (
	"errors"
	"testing"

	"github.com/zahidcakici/champions-league/internal/models"
)

// MockTeamRepository implements repository.TeamRepository for testing
type mockTeamRepository struct {
	teams      []models.Team
	createErr  error
	findAllErr error
	deleteErr  error
	seedErr    error
	seedCalled bool
	deletedIDs []uint
}

func (m *mockTeamRepository) Create(team *models.Team) error {
	if m.createErr != nil {
		return m.createErr
	}
	team.ID = uint(len(m.teams) + 1)
	m.teams = append(m.teams, *team)
	return nil
}

func (m *mockTeamRepository) FindAll() ([]models.Team, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	return m.teams, nil
}

func (m *mockTeamRepository) FindByID(id uint) (*models.Team, error) {
	for _, team := range m.teams {
		if team.ID == id {
			return &team, nil
		}
	}
	return nil, errors.New("team not found")
}

func (m *mockTeamRepository) FindByName(name string) (*models.Team, error) {
	for _, team := range m.teams {
		if team.Name == name {
			return &team, nil
		}
	}
	return nil, errors.New("team not found")
}

func (m *mockTeamRepository) Count() (int64, error) {
	return int64(len(m.teams)), nil
}

func (m *mockTeamRepository) Delete(id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.deletedIDs = append(m.deletedIDs, id)
	// Remove from teams slice
	for i, team := range m.teams {
		if team.ID == id {
			m.teams = append(m.teams[:i], m.teams[i+1:]...)
			break
		}
	}
	return nil
}

func (m *mockTeamRepository) DeleteAll() error {
	m.teams = []models.Team{}
	return nil
}

func (m *mockTeamRepository) SeedDefault() error {
	m.seedCalled = true
	if m.seedErr != nil {
		return m.seedErr
	}
	// Only seed if empty
	if len(m.teams) == 0 {
		m.teams = models.DefaultTeams()
		for i := range m.teams {
			m.teams[i].ID = uint(i + 1)
		}
	}
	return nil
}

func TestTeamService_GetAllTeams(t *testing.T) {
	mockRepo := &mockTeamRepository{}
	service := NewTeamService(mockRepo)

	teams, err := service.GetAllTeams()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should auto-seed default teams
	if !mockRepo.seedCalled {
		t.Error("Expected SeedDefault to be called")
	}

	if len(teams) != 4 {
		t.Errorf("Expected 4 default teams, got %d", len(teams))
	}
}

func TestTeamService_GetAllTeams_SeedError(t *testing.T) {
	mockRepo := &mockTeamRepository{
		seedErr: errors.New("seed failed"),
	}
	service := NewTeamService(mockRepo)

	_, err := service.GetAllTeams()
	if err == nil {
		t.Error("Expected error when seed fails")
	}
}

func TestTeamService_GetAllTeams_FindAllError(t *testing.T) {
	mockRepo := &mockTeamRepository{
		findAllErr: errors.New("database error"),
	}
	service := NewTeamService(mockRepo)

	_, err := service.GetAllTeams()
	if err == nil {
		t.Error("Expected error when FindAll fails")
	}
}

func TestTeamService_CreateTeam(t *testing.T) {
	mockRepo := &mockTeamRepository{}
	service := NewTeamService(mockRepo)

	err := service.CreateTeam("New Team", 75)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockRepo.teams) != 1 {
		t.Errorf("Expected 1 team, got %d", len(mockRepo.teams))
	}

	if mockRepo.teams[0].Name != "New Team" {
		t.Errorf("Expected team name 'New Team', got '%s'", mockRepo.teams[0].Name)
	}

	if mockRepo.teams[0].Power != 75 {
		t.Errorf("Expected power 75, got %d", mockRepo.teams[0].Power)
	}
}

func TestTeamService_CreateTeam_Error(t *testing.T) {
	mockRepo := &mockTeamRepository{
		createErr: errors.New("create failed"),
	}
	service := NewTeamService(mockRepo)

	err := service.CreateTeam("New Team", 75)
	if err == nil {
		t.Error("Expected error when create fails")
	}
}

func TestTeamService_DeleteTeam(t *testing.T) {
	mockRepo := &mockTeamRepository{
		teams: []models.Team{
			{ID: 1, Name: "Team A", Power: 80},
			{ID: 2, Name: "Team B", Power: 75},
		},
	}
	service := NewTeamService(mockRepo)

	err := service.DeleteTeam(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockRepo.deletedIDs) != 1 || mockRepo.deletedIDs[0] != 1 {
		t.Error("Expected team ID 1 to be deleted")
	}

	if len(mockRepo.teams) != 1 {
		t.Errorf("Expected 1 team remaining, got %d", len(mockRepo.teams))
	}
}

func TestTeamService_DeleteTeam_Error(t *testing.T) {
	mockRepo := &mockTeamRepository{
		deleteErr: errors.New("delete failed"),
	}
	service := NewTeamService(mockRepo)

	err := service.DeleteTeam(1)
	if err == nil {
		t.Error("Expected error when delete fails")
	}
}

func TestTeamService_SeedTeams(t *testing.T) {
	mockRepo := &mockTeamRepository{}
	service := NewTeamService(mockRepo)

	err := service.SeedTeams()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !mockRepo.seedCalled {
		t.Error("Expected SeedDefault to be called")
	}
}

func TestTeamService_CreateMultipleTeams(t *testing.T) {
	mockRepo := &mockTeamRepository{}
	service := NewTeamService(mockRepo)

	teamsToCreate := []struct {
		name  string
		power int
	}{
		{"Team Alpha", 90},
		{"Team Beta", 85},
		{"Team Gamma", 70},
	}

	for _, tc := range teamsToCreate {
		err := service.CreateTeam(tc.name, tc.power)
		if err != nil {
			t.Fatalf("Failed to create team %s: %v", tc.name, err)
		}
	}

	if len(mockRepo.teams) != 3 {
		t.Errorf("Expected 3 teams, got %d", len(mockRepo.teams))
	}

	// Verify each team has unique ID
	ids := make(map[uint]bool)
	for _, team := range mockRepo.teams {
		if ids[team.ID] {
			t.Errorf("Duplicate team ID: %d", team.ID)
		}
		ids[team.ID] = true
	}
}

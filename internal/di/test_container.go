package di

import (
	"conduit/internal/repository/postgres"
	"conduit/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Mock StrategyService
type MockStrategyService struct {
	mock.Mock
}

func TestBuildContainer(t *testing.T) {
	// Override the NewDB function
	origNewDB := postgres.NewDB
	postgres.NewDB = func() (*postgres.DB, error) {
		return &postgres.DB{}, nil
	}
	defer func() { postgres.NewDB = origNewDB }()

	// Override the NewStrategyService function
	origNewStrategyService := services.NewStrategyService
	services.NewStrategyService = func(db *postgres.DB) *services.StrategyService {
		return &services.StrategyService{}
	}
	defer func() { services.NewStrategyService = origNewStrategyService }()

	// Test
	container, err := BuildContainer()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, container)

	// Test getting services from the container
	db, err := container.SafeGet("db")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	strategyService, err := container.SafeGet("strategy_service")
	assert.NoError(t, err)
	assert.NotNil(t, strategyService)
}
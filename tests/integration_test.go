package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/internal/repository/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This integration test requires a running DB.
// It will only run if DATABASE_URL is set, otherwise it skips.
// Use: export DATABASE_URL=... && go test ./tests/...

func TestIntegration_UserFlow(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		t.Skip("Skipping integration test: DATABASE_URL not set")
	}

	// 1. Setup DB Connection
	// 1. Setup DB Connection
	config, err := pgxpool.ParseConfig(connStr)
	require.NoError(t, err)

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		registerTypes(t, conn)
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	require.NoError(t, err)
	defer pool.Close()

	// 2. Setup Repo
	repo := postgres.NewUserRepo(pool)

	// 3. Create User
	email := "integration-" + time.Now().Format("20060102-150405") + "@test.com"
	newUser := &domain.User{
		Email:       email,
		Name:        strPtr("Integration Tester"),
		Address:     strPtr("123 Test Lane"),
		Role:        "admin",
		Phone:       strPtr("+1-555-0000"),
		AppMetadata: map[string]interface{}{"source": "test"},
	}

	err = repo.Create(context.Background(), newUser)
	require.NoError(t, err)
	assert.NotEmpty(t, newUser.ID)
	assert.NotZero(t, newUser.CreatedAt)

	// 4. Get User
	fetchedUser, err := repo.GetByID(context.Background(), newUser.ID)
	require.NoError(t, err)
	assert.Equal(t, newUser.Email, fetchedUser.Email)
	assert.NotNil(t, fetchedUser.Name)
	assert.Equal(t, "Integration Tester", *fetchedUser.Name)

	// 5. Cleanup
	_ = repo.Delete(context.Background(), newUser.ID)
}

func strPtr(s string) *string {
	return &s
}

func registerTypes(t *testing.T, conn *pgx.Conn) {
	types := []string{
		"users.user_request",
		"devices.device_request",
		"cabinets.cabinet_request",
	}
	for _, typeName := range types {
		dt, err := conn.LoadType(context.Background(), typeName)
		if err != nil {
			t.Logf("Warning: Could not load type %s: %v", typeName, err)
			continue
		}
		conn.TypeMap().RegisterType(dt)
	}
}

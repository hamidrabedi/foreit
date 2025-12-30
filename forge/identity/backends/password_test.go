package backends

import (
	"context"
	"database/sql"
	"testing"

	"github.com/forgego/forge/db"
	"github.com/forgego/forge/identity/models"
	"github.com/forgego/forge/identity/repository"
	"github.com/forgego/forge/identity/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func setupPasswordBackendTest(t *testing.T) (AuthenticationBackend, *db.DB, context.Context) {
	testDB := setupTestDB(t)
	userRepo := repository.NewUserRepository(testDB)
	backend := NewPasswordBackend(userRepo)
	ctx := context.Background()
	return backend, testDB, ctx
}

func setupTestDB(t *testing.T) *db.DB {
	sqlDB, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	// Ensure single connection for in-memory DB
	sqlDB.SetMaxOpenConns(1)

	testDB := &db.DB{DB: sqlDB, Driver: "sqlite3"}

	_, err = testDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(150) UNIQUE NOT NULL,
			email VARCHAR(254) UNIQUE NOT NULL,
			password VARCHAR(128) NOT NULL,
			first_name VARCHAR(150),
			last_name VARCHAR(150),
			is_active BOOLEAN DEFAULT 1,
			is_staff BOOLEAN DEFAULT 0,
			is_superuser BOOLEAN DEFAULT 0,
			is_locked BOOLEAN DEFAULT 0,
			email_verified BOOLEAN DEFAULT 0,
			phone_number VARCHAR(20),
			phone_verified BOOLEAN DEFAULT 0,
			timezone VARCHAR(50),
			locale VARCHAR(10),
			language VARCHAR(10),
			bio TEXT,
			website VARCHAR(255),
			location VARCHAR(255),
			avatar VARCHAR(255),
			password_changed_at TIMESTAMP,
			password_expires_at TIMESTAMP,
			must_change_password BOOLEAN DEFAULT 0,
			locked_at TIMESTAMP,
			locked_reason VARCHAR(255),
			failed_login_attempts INTEGER DEFAULT 0,
			last_failed_login_at TIMESTAMP,
			email_verified_at TIMESTAMP,
			date_joined TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_login TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP
		);
	`)
	require.NoError(t, err)

	return testDB
}

func TestPasswordBackend_Authenticate(t *testing.T) {
	backend, testDB, ctx := setupPasswordBackendTest(t)
	defer testDB.Close()

	// Create test user
	userRepo := repository.NewUserRepository(testDB)
	hashedPassword, _ := utils.HashPassword("correctpassword")
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
		IsActive: true,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	t.Run("authenticates with correct username and password", func(t *testing.T) {
		credentials := map[string]string{
			"username": "testuser",
			"password": "correctpassword",
		}

		authenticated, err := backend.Authenticate(ctx, credentials)
		require.NoError(t, err)
		assert.Equal(t, user.ID, authenticated.ID)
		assert.Equal(t, user.Username, authenticated.Username)
	})

	t.Run("authenticates with email and password", func(t *testing.T) {
		credentials := map[string]string{
			"email":    "test@example.com",
			"password": "correctpassword",
		}

		authenticated, err := backend.Authenticate(ctx, credentials)
		require.NoError(t, err)
		assert.Equal(t, user.ID, authenticated.ID)
	})

	t.Run("fails with incorrect password", func(t *testing.T) {
		credentials := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}

		_, err := backend.Authenticate(ctx, credentials)
		assert.Error(t, err)
	})

	t.Run("fails with non-existent user", func(t *testing.T) {
		credentials := map[string]string{
			"username": "nonexistent",
			"password": "password123",
		}

		_, err := backend.Authenticate(ctx, credentials)
		assert.Error(t, err)
	})

	t.Run("fails with inactive user", func(t *testing.T) {
		// Deactivate user
		user.IsActive = false
		userRepo.Update(ctx, user)

		credentials := map[string]string{
			"username": "testuser",
			"password": "correctpassword",
		}

		_, err := backend.Authenticate(ctx, credentials)
		assert.Error(t, err)
	})

	t.Run("fails with locked user", func(t *testing.T) {
		// Reactivate and lock user
		user.IsActive = true
		user.IsLocked = true
		userRepo.Update(ctx, user)

		credentials := map[string]string{
			"username": "testuser",
			"password": "correctpassword",
		}

		_, err := backend.Authenticate(ctx, credentials)
		assert.Error(t, err)
	})
}

func TestPasswordBackend_Supports(t *testing.T) {
	backend, testDB, _ := setupPasswordBackendTest(t)
	defer testDB.Close()

	t.Run("supports password credential type", func(t *testing.T) {
		assert.True(t, backend.Supports("password"))
	})

	t.Run("does not support other credential types", func(t *testing.T) {
		assert.False(t, backend.Supports("token"))
		assert.False(t, backend.Supports("oauth"))
	})
}

func TestPasswordBackend_Name(t *testing.T) {
	backend, testDB, _ := setupPasswordBackendTest(t)
	defer testDB.Close()

	assert.Equal(t, "password", backend.Name())
}

func TestPasswordBackend_GetUser(t *testing.T) {
	backend, testDB, ctx := setupPasswordBackendTest(t)
	defer testDB.Close()

	t.Run("returns nil for password backend", func(t *testing.T) {
		// Password backend doesn't support GetUser by identifier
		// It returns nil, nil to indicate it's not applicable
		user, err := backend.GetUser(ctx, "identifier")
		assert.NoError(t, err)
		assert.Nil(t, user)
	})
}

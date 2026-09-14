package generate

import (
	"testing"

	"github.com/forgego/forge/db/migrate/core"
	"github.com/forgego/forge/db/migrate/sql"
	"github.com/forgego/forge/db/migrate/state"
	"github.com/stretchr/testify/require"
)

func TestNewMigrationGeneratorPreservesCallerDriver(t *testing.T) {
	driver := core.DriverSQLite
	builder, err := sql.NewSQLBuilder(driver)
	require.NoError(t, err)

	gen, err := NewMigrationGenerator("models", "migrations", driver, NewDetector(), builder, state.NewInMemoryState())

	require.NoError(t, err)
	require.Equal(t, driver, gen.driver)
}

func TestNewMigrationGeneratorWithDefaultsUsesCallerDriver(t *testing.T) {
	gen, err := NewMigrationGeneratorWithDefaults("models", "migrations", core.DriverSQLite)

	require.NoError(t, err)
	require.Equal(t, core.DriverSQLite, gen.driver)
}

func TestNewMigrationGeneratorWithDefaultsRejectsUnsupportedCallerDriver(t *testing.T) {
	gen, err := NewMigrationGeneratorWithDefaults("models", "migrations", core.Driver("unsupported"))

	require.Nil(t, gen)
	require.ErrorContains(t, err, "unsupported database driver: unsupported")
}

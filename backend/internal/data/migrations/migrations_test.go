// Package migrations
package migrations

import (
	"database/sql"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sysadminsmedia/homebox/backend/internal/sys/config"

	_ "github.com/sysadminsmedia/homebox/backend/internal/data/migrations/sqlite3"
	_ "modernc.org/sqlite"
)

// TestSQLiteMigrationsProduceBuildSchema runs the full sqlite3 migration chain
// against a fresh database and asserts the resulting schema matches what the
// ent-generated queries expect. This guards against partial renames: the
// setups→builds rename originally renamed the tables but left the
// group_setups FK column behind, breaking every build query at runtime.
func TestSQLiteMigrationsProduceBuildSchema(t *testing.T) {
	db, err := sql.Open("sqlite", t.TempDir()+"/homebox.db")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	fs, err := Migrations(config.DriverSqlite3)
	require.NoError(t, err)

	goose.SetBaseFS(fs)
	require.NoError(t, goose.SetDialect("sqlite3"))
	require.NoError(t, goose.Up(db, "sqlite3"))

	// The renamed tables must exist with the renamed FK column.
	var column string
	require.NoError(t, db.QueryRow(
		`SELECT name FROM pragma_table_info('builds') WHERE name = 'group_builds'`,
	).Scan(&column))
	assert.Equal(t, "group_builds", column)

	require.NoError(t, db.QueryRow(
		`SELECT name FROM pragma_table_info('build_parts') WHERE name = 'build_id'`,
	).Scan(&column))
	assert.Equal(t, "build_id", column)

	// No pre-rename tables or columns may survive.
	for _, leftover := range []struct{ table, column string }{
		{"setups", ""},
		{"setup_parts", ""},
		{"builds", "group_setups"},
		{"build_parts", "setup_id"},
	} {
		var count int
		if leftover.column == "" {
			require.NoError(t, db.QueryRow(
				`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`,
				leftover.table,
			).Scan(&count))
		} else {
			require.NoError(t, db.QueryRow(
				`SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = ?`,
				leftover.table, leftover.column,
			).Scan(&count))
		}
		assert.Zero(t, count, "leftover %s.%s should have been renamed", leftover.table, leftover.column)
	}
}

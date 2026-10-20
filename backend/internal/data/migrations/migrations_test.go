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

// TestSQLiteMigrationEntitySportBackfill runs the migration chain up to just
// before the entity sport migration, seeds legacy rows (prefixed template
// names, bmx builds), then applies it and asserts the backfill: sport is
// derived from the name prefix, prefixes are stripped, and bmx builds are
// migrated to cycling.
func TestSQLiteMigrationEntitySportBackfill(t *testing.T) {
	db, err := sql.Open("sqlite", t.TempDir()+"/homebox.db")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	fs, err := Migrations(config.DriverSqlite3)
	require.NoError(t, err)

	goose.SetBaseFS(fs)
	require.NoError(t, goose.SetDialect("sqlite3"))
	require.NoError(t, goose.UpTo(db, "sqlite3", 20261022000000))

	_, err = db.Exec(`
		INSERT INTO groups (id, created_at, updated_at, name, currency) VALUES
		('00000000-0000-0000-0000-000000000001', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'g', 'USD');
		INSERT INTO entity_types (id, created_at, updated_at, name, is_location, group_entity_types) VALUES
		('00000000-0000-0000-0000-0000000000e1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Item', false, '00000000-0000-0000-0000-000000000001');
		INSERT INTO entities (id, created_at, updated_at, group_entities, name, quantity, entity_type_entities) VALUES
		('00000000-0000-0000-0000-0000000000aa', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '00000000-0000-0000-0000-000000000001', 'Skateboard Deck', 1, '00000000-0000-0000-0000-0000000000e1'),
		('00000000-0000-0000-0000-0000000000ab', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '00000000-0000-0000-0000-000000000001', 'Bike Chain', 1, '00000000-0000-0000-0000-0000000000e1'),
		('00000000-0000-0000-0000-0000000000ac', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '00000000-0000-0000-0000-000000000001', 'Generic Box', 1, '00000000-0000-0000-0000-0000000000e1');
		INSERT INTO builds (id, created_at, updated_at, group_builds, name, sport) VALUES
		('00000000-0000-0000-0000-0000000000ba', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '00000000-0000-0000-0000-000000000001', 'Old Ride', 'bmx');
	`)
	require.NoError(t, err)

	require.NoError(t, goose.Up(db, "sqlite3"))

	type row struct{ name, sport string }
	rows := map[string]row{}
	r, err := db.Query(`SELECT name, COALESCE(sport, '') FROM entities WHERE name IN ('Deck', 'Chain', 'Generic Box')`)
	require.NoError(t, err)
	for r.Next() {
		var e row
		require.NoError(t, r.Scan(&e.name, &e.sport))
		rows[e.name] = e
	}
	require.NoError(t, r.Close())
	assert.Equal(t, row{"Deck", "skateboarding"}, rows["Deck"])
	assert.Equal(t, row{"Chain", "cycling"}, rows["Chain"])
	// Generic entities keep their name and stay sportless — never dropped.
	assert.Equal(t, row{"Generic Box", ""}, rows["Generic Box"])

	var sport string
	require.NoError(t, db.QueryRow(
		`SELECT sport FROM builds WHERE id = '00000000-0000-0000-0000-0000000000ba'`,
	).Scan(&sport))
	assert.Equal(t, "cycling", sport)
}

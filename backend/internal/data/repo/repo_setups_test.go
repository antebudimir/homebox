package repo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupRepository_CRUD(t *testing.T) {
	ctx := context.Background()

	created, err := tRepos.Setups.Create(ctx, tGroup.ID, SetupCreate{
		Name:        "Park Deck",
		Description: "Main street setup",
		Sport:       "skateboarding",
		IsActive:    true,
	})
	require.NoError(t, err)
	assert.NotEqual(t, [16]byte{}, created.ID)
	assert.Equal(t, "Park Deck", created.Name)
	assert.True(t, created.IsActive)

	got, err := tRepos.Setups.GetOne(ctx, tGroup.ID, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.Name, got.Name)

	updated, err := tRepos.Setups.Update(ctx, tGroup.ID, created.ID, SetupUpdate{
		Name:        "Park Deck v2",
		Description: "Updated",
		Sport:       "skateboarding",
		IsActive:    false,
	})
	require.NoError(t, err)
	assert.Equal(t, "Park Deck v2", updated.Name)
	assert.False(t, updated.IsActive)

	all, err := tRepos.Setups.GetAll(ctx, tGroup.ID)
	require.NoError(t, err)
	assert.Len(t, all, 1)

	require.NoError(t, tRepos.Setups.Delete(ctx, tGroup.ID, created.ID))
	_, err = tRepos.Setups.GetOne(ctx, tGroup.ID, created.ID)
	assert.Error(t, err)
}

func TestSetupRepository_Parts(t *testing.T) {
	ctx := context.Background()
	entities := useEntities(t, 2)

	setup, err := tRepos.Setups.Create(ctx, tGroup.ID, SetupCreate{Name: "Cruiser"})
	require.NoError(t, err)

	part, err := tRepos.Setups.AddPart(ctx, tGroup.ID, setup.ID, SetupPartCreate{
		EntityID: entities[0].ID,
		Position: "front_truck",
	})
	require.NoError(t, err)
	assert.Equal(t, entities[0].ID, part.EntityID)
	assert.Equal(t, "front_truck", part.Position)
	assert.Nil(t, part.RemovedDate)
	assert.Equal(t, entities[0].Name, part.EntityName)

	// A part from another group must be rejected.
	assert.NotEqual(t, part.EntityID, entities[1].ID)

	current, err := tRepos.Setups.GetParts(ctx, tGroup.ID, setup.ID, true)
	require.NoError(t, err)
	assert.Len(t, current, 1)

	// Removing stamps removed_date rather than deleting the row.
	require.NoError(t, tRepos.Setups.RemovePart(ctx, tGroup.ID, part.ID))
	current, err = tRepos.Setups.GetParts(ctx, tGroup.ID, setup.ID, true)
	require.NoError(t, err)
	assert.Len(t, current, 0)

	history, err := tRepos.Setups.GetParts(ctx, tGroup.ID, setup.ID, false)
	require.NoError(t, err)
	assert.Len(t, history, 1)
	require.NotNil(t, history[0].RemovedDate)

	byEntity, err := tRepos.Setups.GetPartsByEntity(ctx, tGroup.ID, entities[0].ID)
	require.NoError(t, err)
	assert.Len(t, byEntity, 1)

	// Removing twice fails: the row is already stamped.
	assert.Error(t, tRepos.Setups.RemovePart(ctx, tGroup.ID, part.ID))
	require.NoError(t, tRepos.Setups.DeletePart(ctx, tGroup.ID, part.ID))

	history, err = tRepos.Setups.GetParts(ctx, tGroup.ID, setup.ID, false)
	require.NoError(t, err)
	assert.Len(t, history, 0)

	// Cleanup
	require.NoError(t, tRepos.Setups.Delete(ctx, tGroup.ID, setup.ID))
}

func TestSetupRepository_Sessions(t *testing.T) {
	ctx := context.Background()

	setup, err := tRepos.Setups.Create(ctx, tGroup.ID, SetupCreate{Name: "BMX"})
	require.NoError(t, err)

	session, err := tRepos.Setups.CreateSession(ctx, tGroup.ID, setup.ID, UsageSessionCreate{
		Date:            time.Now(),
		DurationMinutes: 90,
		Location:        "skatepark",
		Conditions:      "dry",
		Notes:           "first sesh",
		Rating:          4,
	})
	require.NoError(t, err)
	assert.Equal(t, 90, session.DurationMinutes)
	assert.Equal(t, 4, session.Rating)

	sessions, err := tRepos.Setups.GetSessions(ctx, tGroup.ID, setup.ID)
	require.NoError(t, err)
	assert.Len(t, sessions, 1)

	updated, err := tRepos.Setups.UpdateSession(ctx, tGroup.ID, session.ID, UsageSessionUpdate{
		Date:            time.Now(),
		DurationMinutes: 120,
		Location:        "skatepark",
		Conditions:      "hot",
		Notes:           "longer sesh",
		Rating:          5,
	})
	require.NoError(t, err)
	assert.Equal(t, 120, updated.DurationMinutes)
	assert.Equal(t, 5, updated.Rating)

	require.NoError(t, tRepos.Setups.DeleteSession(ctx, tGroup.ID, session.ID))
	sessions, err = tRepos.Setups.GetSessions(ctx, tGroup.ID, setup.ID)
	require.NoError(t, err)
	assert.Len(t, sessions, 0)

	require.NoError(t, tRepos.Setups.Delete(ctx, tGroup.ID, setup.ID))
}
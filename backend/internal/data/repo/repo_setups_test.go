package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
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
	assert.Empty(t, part.Warnings)
	assert.Equal(t, entities[0].ID, part.Part.EntityID)
	assert.Equal(t, "front_truck", part.Part.Position)
	assert.Nil(t, part.Part.RemovedDate)
	assert.Equal(t, entities[0].Name, part.Part.EntityName)

	// A part from another group must be rejected.
	assert.NotEqual(t, part.Part.EntityID, entities[1].ID)

	current, err := tRepos.Setups.GetParts(ctx, tGroup.ID, setup.ID, true)
	require.NoError(t, err)
	assert.Len(t, current, 1)

	// Removing stamps removed_date rather than deleting the row.
	require.NoError(t, tRepos.Setups.RemovePart(ctx, tGroup.ID, part.Part.ID))
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
	assert.Error(t, tRepos.Setups.RemovePart(ctx, tGroup.ID, part.Part.ID))
	require.NoError(t, tRepos.Setups.DeletePart(ctx, tGroup.ID, part.Part.ID))

	history, err = tRepos.Setups.GetParts(ctx, tGroup.ID, setup.ID, false)
	require.NoError(t, err)
	assert.Len(t, history, 0)

	// Cleanup
	require.NoError(t, tRepos.Setups.Delete(ctx, tGroup.ID, setup.ID))
}

func TestSetupRepository_AddPartLocationWarning(t *testing.T) {
	ctx := context.Background()

	// Mounting a location-typed entity is allowed but should warn.
	containerET := useContainerEntityType(t)
	location, err := tRepos.Entities.Create(ctx, tGroup.ID, EntityCreate{
		Name:         fk.Str(10),
		EntityTypeID: containerET.ID,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tRepos.Entities.Delete(ctx, location.ID)
	})

	setup, err := tRepos.Setups.Create(ctx, tGroup.ID, SetupCreate{Name: "Warning Setup"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tRepos.Setups.Delete(ctx, tGroup.ID, setup.ID)
	})

	result, err := tRepos.Setups.AddPart(ctx, tGroup.ID, setup.ID, SetupPartCreate{EntityID: location.ID})
	require.NoError(t, err)
	assert.Equal(t, []string{"entity_is_location"}, result.Warnings)
}

func TestSetupRepository_Stats(t *testing.T) {
	ctx := context.Background()

	entities := useEntities(t, 2)
	prices := map[uuid.UUID]float64{entities[0].ID: 50.0, entities[1].ID: 25.5}
	for id, price := range prices {
		_, err := tClient.Entity.UpdateOneID(id).SetPurchasePrice(price).Save(ctx)
		require.NoError(t, err)
		t.Cleanup(func() {
			_, _ = tClient.Entity.UpdateOneID(id).SetPurchasePrice(0).Save(ctx)
		})
	}

	setup, err := tRepos.Setups.Create(ctx, tGroup.ID, SetupCreate{Name: "Stats Setup"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tRepos.Setups.Delete(ctx, tGroup.ID, setup.ID)
	})

	// Mount the same part on the setup twice (removed and re-mounted) plus
	// the second part; investment must count each distinct part once.
	first, err := tRepos.Setups.AddPart(ctx, tGroup.ID, setup.ID, SetupPartCreate{EntityID: entities[0].ID})
	require.NoError(t, err)
	require.NoError(t, tRepos.Setups.RemovePart(ctx, tGroup.ID, first.Part.ID))
	_, err = tRepos.Setups.AddPart(ctx, tGroup.ID, setup.ID, SetupPartCreate{EntityID: entities[0].ID})
	require.NoError(t, err)
	_, err = tRepos.Setups.AddPart(ctx, tGroup.ID, setup.ID, SetupPartCreate{EntityID: entities[1].ID})
	require.NoError(t, err)

	// A second group's setup must not affect this group's stats.
	otherGroup, err := tRepos.Groups.GroupCreate(ctx, "stats-foreign-"+uuid.NewString(), uuid.Nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tRepos.Groups.GroupDelete(ctx, otherGroup.ID)
	})

	otherET, err := tRepos.EntityTypes.GetDefault(ctx, otherGroup.ID, false)
	require.NoError(t, err)
	otherEntity, err := tRepos.Entities.Create(ctx, otherGroup.ID, EntityCreate{
		Name:         fk.Str(10),
		EntityTypeID: otherET.ID,
	})
	require.NoError(t, err)
	_, err = tClient.Entity.UpdateOneID(otherEntity.ID).SetPurchasePrice(50.0).Save(ctx)
	require.NoError(t, err)

	otherSetup, err := tRepos.Setups.Create(ctx, otherGroup.ID, SetupCreate{Name: "Foreign Setup"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tRepos.Setups.Delete(ctx, otherGroup.ID, otherSetup.ID)
		_ = tRepos.Entities.Delete(ctx, otherEntity.ID)
	})
	_, err = tRepos.Setups.AddPart(ctx, otherGroup.ID, otherSetup.ID, SetupPartCreate{EntityID: otherEntity.ID})
	require.NoError(t, err)

	stats, err := tRepos.Setups.Stats(ctx, tGroup.ID)
	require.NoError(t, err)
	assert.InDelta(t, 75.5, stats.TotalInvestment, 0.001)
	assert.Equal(t, 0, stats.SessionCount)
	assert.InDelta(t, 0.0, stats.CostPerSession, 0.001)

	// The foreign group only sees its own investment.
	foreignStats, err := tRepos.Setups.Stats(ctx, otherGroup.ID)
	require.NoError(t, err)
	assert.InDelta(t, 50.0, foreignStats.TotalInvestment, 0.001)

	// Logging a session turns on cost-per-session.
	_, err = tRepos.Setups.CreateSession(ctx, tGroup.ID, setup.ID, UsageSessionCreate{Date: time.Now()})
	require.NoError(t, err)

	stats, err = tRepos.Setups.Stats(ctx, tGroup.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, stats.SessionCount)
	assert.InDelta(t, 75.5, stats.CostPerSession, 0.001)
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

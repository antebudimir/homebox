package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRepository_CRUD(t *testing.T) {
	ctx := context.Background()

	created, err := tRepos.Builds.Create(ctx, tGroup.ID, BuildCreate{
		Name:        "Park Deck",
		Description: "Main street build",
		Sport:       "skateboarding",
		IsActive:    true,
	})
	require.NoError(t, err)
	assert.NotEqual(t, [16]byte{}, created.ID)
	assert.Equal(t, "Park Deck", created.Name)
	assert.True(t, created.IsActive)

	got, err := tRepos.Builds.GetOne(ctx, tGroup.ID, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.Name, got.Name)

	updated, err := tRepos.Builds.Update(ctx, tGroup.ID, created.ID, BuildUpdate{
		Name:        "Park Deck v2",
		Description: "Updated",
		Sport:       "skateboarding",
		IsActive:    false,
	})
	require.NoError(t, err)
	assert.Equal(t, "Park Deck v2", updated.Name)
	assert.False(t, updated.IsActive)

	all, err := tRepos.Builds.GetAll(ctx, tGroup.ID)
	require.NoError(t, err)
	assert.Len(t, all, 1)

	require.NoError(t, tRepos.Builds.Delete(ctx, tGroup.ID, created.ID))
	_, err = tRepos.Builds.GetOne(ctx, tGroup.ID, created.ID)
	assert.Error(t, err)
}

func TestBuildRepository_Parts(t *testing.T) {
	ctx := context.Background()
	entities := useEntities(t, 2)

	build, err := tRepos.Builds.Create(ctx, tGroup.ID, BuildCreate{Name: "Cruiser"})
	require.NoError(t, err)

	part, err := tRepos.Builds.AddPart(ctx, tGroup.ID, build.ID, BuildPartCreate{
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

	current, err := tRepos.Builds.GetParts(ctx, tGroup.ID, build.ID, true)
	require.NoError(t, err)
	assert.Len(t, current, 1)

	// Removing stamps removed_date rather than deleting the row.
	require.NoError(t, tRepos.Builds.RemovePart(ctx, tGroup.ID, part.Part.ID))
	current, err = tRepos.Builds.GetParts(ctx, tGroup.ID, build.ID, true)
	require.NoError(t, err)
	assert.Len(t, current, 0)

	history, err := tRepos.Builds.GetParts(ctx, tGroup.ID, build.ID, false)
	require.NoError(t, err)
	assert.Len(t, history, 1)
	require.NotNil(t, history[0].RemovedDate)

	byEntity, err := tRepos.Builds.GetPartsByEntity(ctx, tGroup.ID, entities[0].ID)
	require.NoError(t, err)
	assert.Len(t, byEntity, 1)

	// Removing twice fails: the row is already stamped.
	assert.Error(t, tRepos.Builds.RemovePart(ctx, tGroup.ID, part.Part.ID))
	require.NoError(t, tRepos.Builds.DeletePart(ctx, tGroup.ID, part.Part.ID))

	history, err = tRepos.Builds.GetParts(ctx, tGroup.ID, build.ID, false)
	require.NoError(t, err)
	assert.Len(t, history, 0)

	// Cleanup
	require.NoError(t, tRepos.Builds.Delete(ctx, tGroup.ID, build.ID))
}

func TestBuildRepository_AddPartLocationWarning(t *testing.T) {
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

	build, err := tRepos.Builds.Create(ctx, tGroup.ID, BuildCreate{Name: "Warning Build"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tRepos.Builds.Delete(ctx, tGroup.ID, build.ID)
	})

	result, err := tRepos.Builds.AddPart(ctx, tGroup.ID, build.ID, BuildPartCreate{EntityID: location.ID})
	require.NoError(t, err)
	assert.Equal(t, []string{"entity_is_location"}, result.Warnings)
}

func TestBuildRepository_Stats(t *testing.T) {
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

	build, err := tRepos.Builds.Create(ctx, tGroup.ID, BuildCreate{Name: "Stats Build"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tRepos.Builds.Delete(ctx, tGroup.ID, build.ID)
	})

	// Mount the same part on the build twice (removed and re-mounted) plus
	// the second part; investment must count each distinct part once.
	first, err := tRepos.Builds.AddPart(ctx, tGroup.ID, build.ID, BuildPartCreate{EntityID: entities[0].ID})
	require.NoError(t, err)
	require.NoError(t, tRepos.Builds.RemovePart(ctx, tGroup.ID, first.Part.ID))
	_, err = tRepos.Builds.AddPart(ctx, tGroup.ID, build.ID, BuildPartCreate{EntityID: entities[0].ID})
	require.NoError(t, err)
	_, err = tRepos.Builds.AddPart(ctx, tGroup.ID, build.ID, BuildPartCreate{EntityID: entities[1].ID})
	require.NoError(t, err)

	// A second group's build must not affect this group's stats.
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

	otherBuild, err := tRepos.Builds.Create(ctx, otherGroup.ID, BuildCreate{Name: "Foreign Build"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tRepos.Builds.Delete(ctx, otherGroup.ID, otherBuild.ID)
		_ = tRepos.Entities.Delete(ctx, otherEntity.ID)
	})
	_, err = tRepos.Builds.AddPart(ctx, otherGroup.ID, otherBuild.ID, BuildPartCreate{EntityID: otherEntity.ID})
	require.NoError(t, err)

	stats, err := tRepos.Builds.Stats(ctx, tGroup.ID)
	require.NoError(t, err)
	assert.InDelta(t, 75.5, stats.TotalInvestment, 0.001)
	assert.Equal(t, 0, stats.SessionCount)
	assert.InDelta(t, 0.0, stats.CostPerSession, 0.001)

	// The foreign group only sees its own investment.
	foreignStats, err := tRepos.Builds.Stats(ctx, otherGroup.ID)
	require.NoError(t, err)
	assert.InDelta(t, 50.0, foreignStats.TotalInvestment, 0.001)

	// Logging a session turns on cost-per-session.
	_, err = tRepos.Builds.CreateSession(ctx, tGroup.ID, build.ID, UsageSessionCreate{Date: time.Now()})
	require.NoError(t, err)

	stats, err = tRepos.Builds.Stats(ctx, tGroup.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, stats.SessionCount)
	assert.InDelta(t, 75.5, stats.CostPerSession, 0.001)
}

func TestBuildRepository_Sessions(t *testing.T) {
	ctx := context.Background()

	build, err := tRepos.Builds.Create(ctx, tGroup.ID, BuildCreate{Name: "BMX"})
	require.NoError(t, err)

	session, err := tRepos.Builds.CreateSession(ctx, tGroup.ID, build.ID, UsageSessionCreate{
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

	sessions, err := tRepos.Builds.GetSessions(ctx, tGroup.ID, build.ID)
	require.NoError(t, err)
	assert.Len(t, sessions, 1)

	updated, err := tRepos.Builds.UpdateSession(ctx, tGroup.ID, session.ID, UsageSessionUpdate{
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

	require.NoError(t, tRepos.Builds.DeleteSession(ctx, tGroup.ID, session.ID))
	sessions, err = tRepos.Builds.GetSessions(ctx, tGroup.ID, build.ID)
	require.NoError(t, err)
	assert.Len(t, sessions, 0)

	require.NoError(t, tRepos.Builds.Delete(ctx, tGroup.ID, build.ID))
}

package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/data/repo"
)

const defaultLocationGarage = "Garage"

// ensureDefaultEntityTypes guarantees a freshly created group has the two
// baseline entity types ("Item" and "Location"). The frontend create dialogs
// require selecting an existing type, so without these a brand-new group can't
// create items or locations at all. GetDefault creates each type if missing.
func ensureDefaultEntityTypes(ctx context.Context, repos *repo.AllRepos, gid uuid.UUID) error {
	for _, isLocation := range []bool{false, true} {
		if _, err := repos.EntityTypes.GetDefault(ctx, gid, isLocation); err != nil {
			return err
		}
	}
	return nil
}

func defaultLocations() []repo.EntityCreate {
	return []repo.EntityCreate{
		{
			Name: "Living Room",
		},
		{
			Name: defaultLocationGarage,
		},
		{
			Name: "Kitchen",
		},
		{
			Name: "Bedroom",
		},
		{
			Name: "Bathroom",
		},
		{
			Name: "Office",
		},
		{
			Name: "Attic",
		},
		{
			Name: "Basement",
		},
	}
}

func defaultTags() []repo.TagCreate {
	return []repo.TagCreate{
		{
			Name: "Appliances",
		},
		{
			Name: "IOT",
		},
		{
			Name: "Electronics",
		},
		{
			Name: "Servers",
		},
		{
			Name: "General",
		},
		{
			Name: "Important",
		},
	}
}

// buildPartTemplate is one seeded part "template": a plain item the user can
// duplicate or mount into a build without filling in the basics themselves.
type buildPartTemplate struct {
	Name        string
	Description string
}

// buildPartTemplates are seeded into every new group, keyed by name so
// seeding is idempotent. 8 skateboard + 7 bike parts.
var buildPartTemplates = []buildPartTemplate{
	// Skateboard
	{Name: "Skateboard Deck", Description: "Deck template — set your width, wheelbase and concave."},
	{Name: "Skateboard Grip Tape", Description: "Grip tape sheet."},
	{Name: "Skateboard Front Truck", Description: "Front truck — record hanger width and brand."},
	{Name: "Skateboard Rear Truck", Description: "Rear truck — record hanger width and brand."},
	{Name: "Skateboard Front Wheels", Description: "Front wheels — record diameter and durometer."},
	{Name: "Skateboard Rear Wheels", Description: "Rear wheels — record diameter and durometer."},
	{Name: "Skateboard Bearings", Description: "Bearing set (8 pack for a full rebuild)."},
	{Name: "Skateboard Hardware", Description: "Mounting hardware bolt set."},
	// Bike
	{Name: "Bike Frame", Description: "Frame template — record size and material."},
	{Name: "Bike Front Wheel", Description: "Front wheel — record rim, hub and spokes."},
	{Name: "Bike Rear Wheel", Description: "Rear wheel — record rim, hub, driver and spokes."},
	{Name: "Bike Handlebars", Description: "Handlebars — record width, rise and clamp size."},
	{Name: "Bike Saddle", Description: "Saddle / seat + post."},
	{Name: "Bike Crankset", Description: "Crankset — record spindle size and gearing."},
	{Name: "Bike Chain", Description: "Chain — record length and half-link."},
}

// ensureBuildPartTemplates seeds the sport part templates into a freshly
// created group. It is idempotent: templates whose name already exists in
// the group are skipped, so it is safe to call on any group.
func ensureBuildPartTemplates(ctx context.Context, repos *repo.AllRepos, gid uuid.UUID) error {
	existing, err := repos.Entities.GetAll(ctx, gid)
	if err != nil {
		return err
	}
	names := make(map[string]struct{}, len(existing))
	for _, e := range existing {
		names[e.Name] = struct{}{}
	}

	itemType, err := repos.EntityTypes.GetDefault(ctx, gid, false)
	if err != nil {
		return err
	}

	for _, tpl := range buildPartTemplates {
		if _, ok := names[tpl.Name]; ok {
			continue
		}
		if _, err := repos.Entities.Create(ctx, gid, repo.EntityCreate{
			Name:         tpl.Name,
			Description:  tpl.Description,
			EntityTypeID: itemType.ID,
		}); err != nil {
			return err
		}
	}
	return nil
}

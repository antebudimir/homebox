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
	Sport       string
}

// buildPartTemplates are seeded into every new group, keyed by name so
// seeding is idempotent. 9 skateboard + 7 bike parts.
var buildPartTemplates = []buildPartTemplate{
	// Skateboarding
	{Name: "Deck", Description: "Deck template — set your width, wheelbase and concave.", Sport: "skateboarding"},
	{Name: "Grip Tape", Description: "Grip tape sheet.", Sport: "skateboarding"},
	{Name: "Front Truck", Description: "Front truck — record hanger width and brand.", Sport: "skateboarding"},
	{Name: "Rear Truck", Description: "Rear truck — record hanger width and brand.", Sport: "skateboarding"},
	{Name: "Front Wheels", Description: "Front wheels — record diameter and durometer.", Sport: "skateboarding"},
	{Name: "Rear Wheels", Description: "Rear wheels — record diameter and durometer.", Sport: "skateboarding"},
	{Name: "Bearings", Description: "Bearing set (8 pack for a full rebuild).", Sport: "skateboarding"},
	{Name: "Hardware", Description: "4 bolts, 4 nuts, 4 washers.", Sport: "skateboarding"},
	{Name: "Bushings", Description: "4 bushings (or 2 pairs).", Sport: "skateboarding"},
	// Cycling
	{Name: "Frame", Description: "Frame template — record size and material.", Sport: "cycling"},
	{Name: "Front Wheel", Description: "Front wheel — record rim, hub and spokes.", Sport: "cycling"},
	{Name: "Rear Wheel", Description: "Rear wheel — record rim, hub, driver and spokes.", Sport: "cycling"},
	{Name: "Handlebars", Description: "Handlebars — record width, rise and clamp size.", Sport: "cycling"},
	{Name: "Saddle", Description: "Saddle / seat + post.", Sport: "cycling"},
	{Name: "Crankset", Description: "Crankset — record spindle size and gearing.", Sport: "cycling"},
	{Name: "Chain", Description: "Chain — record length and half-link.", Sport: "cycling"},
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
			Sport:        tpl.Sport,
			EntityTypeID: itemType.ID,
		}); err != nil {
			return err
		}
	}
	return nil
}

package repo

// SportPositionConfig is the typed, pre-defined set of mount positions for
// one sport. Positions are code config rather than an entity to keep the
// surface area small; this can graduate to a real entity later if groups
// ever need custom positions.
type SportPositionConfig struct {
	Slug      string   `json:"slug"`
	Label     string   `json:"label"`
	Positions []string `json:"positions"`
}

// SportPositionConfigs maps each supported sport to its mount positions.
var SportPositionConfigs = []SportPositionConfig{
	{
		Slug:  "skateboarding",
		Label: "Skateboarding",
		// Skateboards have four wheels tracked as pairs, so wheels and
		// bearings mount to the truck they live on rather than having
		// their own positions.
		Positions: []string{"deck", "front_truck", "rear_truck"},
	},
	{
		Slug:      "cycling",
		Label:     "Cycling",
		Positions: []string{"frame", "front_wheel", "rear_wheel", "handlebars", "saddle"},
	},
}

// SportPositions returns the mount positions for a sport slug, or nil when
// the sport is unknown (free-text sport values stay valid and unvalidated).
func SportPositions(slug string) []string {
	for _, s := range SportPositionConfigs {
		if s.Slug == slug {
			return s.Positions
		}
	}
	return nil
}

// Sports returns the supported sport/position configs.
func (r *BuildRepository) Sports() []SportPositionConfig {
	return SportPositionConfigs
}

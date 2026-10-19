package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/schema/mixins"
)

// UsageSession logs a single use of a Setup (date, duration, location,
// conditions, rating). Sessions drive cost-per-use and wear analytics.
type UsageSession struct {
	ent.Schema
}

func (UsageSession) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.BaseMixin{},
	}
}

// Fields of the UsageSession.
func (UsageSession) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("setup_id", uuid.UUID{}),
		field.Time("date").
			Optional(),
		field.Int("duration_minutes").
			Default(0),
		field.String("location").
			MaxLen(255).
			Optional(),
		field.String("conditions").
			MaxLen(255).
			Optional(),
		field.String("notes").
			MaxLen(1000).
			Optional(),
		field.Int("rating").
			Optional(),
	}
}

// Edges of the UsageSession.
func (UsageSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("setup", Setup.Type).
			Ref("usage_sessions").
			Field("setup_id").
			Required().
			Unique(),
	}
}
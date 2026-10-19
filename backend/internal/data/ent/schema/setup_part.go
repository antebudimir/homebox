package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/schema/mixins"
)

// SetupPart is the join between a Setup and a part (Entity). A null
// removed_date means the part is currently mounted; a value marks when it was
// taken off, keeping the full mounting history. Position is a free-form label
// (e.g. "front_truck", "rear_wheel") so it works for any sport.
type SetupPart struct {
	ent.Schema
}

func (SetupPart) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.BaseMixin{},
	}
}

// Fields of the SetupPart.
func (SetupPart) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("setup_id", uuid.UUID{}),
		field.UUID("entity_id", uuid.UUID{}),
		field.Time("added_date").
			Default(time.Now),
		field.Time("removed_date").
			Optional(),
		field.String("position").
			MaxLen(64).
			Optional(),
	}
}

// Edges of the SetupPart.
func (SetupPart) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("setup", Setup.Type).
			Ref("parts").
			Field("setup_id").
			Required().
			Unique(),
		edge.From("entity", Entity.Type).
			Ref("setup_parts").
			Field("entity_id").
			Required().
			Unique(),
	}
}
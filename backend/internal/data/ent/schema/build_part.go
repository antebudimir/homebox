package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/schema/mixins"
)

// BuildPart is the join between a Build and a part (Entity). A null
// removed_date means the part is currently mounted; a value marks when it was
// taken off, keeping the full mounting history. Position is a free-form label
// (e.g. "front_truck", "rear_wheel") so it works for any sport.
type BuildPart struct {
	ent.Schema
}

func (BuildPart) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.BaseMixin{},
	}
}

// Fields of the BuildPart.
func (BuildPart) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("build_id", uuid.UUID{}),
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

// Edges of the BuildPart.
func (BuildPart) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("build", Build.Type).
			Ref("parts").
			Field("build_id").
			Required().
			Unique(),
		edge.From("entity", Entity.Type).
			Ref("build_parts").
			Field("entity_id").
			Required().
			Unique(),
	}
}

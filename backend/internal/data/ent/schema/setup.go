package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/schema/mixins"
)

// Setup holds the schema definition for the Setup entity. A Setup is a named
// assembly of parts (e.g. a skateboard or bike build). Parts are attached via
// SetupPart rows which keep added/removed dates, so a part's mounting history
// across setups is preserved.
type Setup struct {
	ent.Schema
}

func (Setup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.BaseMixin{},
		mixins.DetailsMixin{},
		GroupMixin{ref: "setups"},
	}
}

// Fields of the Setup.
func (Setup) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("is_active").
			Default(false),
		field.String("sport").
			MaxLen(64).
			Optional(),
	}
}

// Edges of the Setup.
func (Setup) Edges() []ent.Edge {
	owned := func(s string, t any) ent.Edge {
		return edge.To(s, t).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			})
	}

	return []ent.Edge{
		owned("parts", SetupPart.Type),
		owned("usage_sessions", UsageSession.Type),
	}
}
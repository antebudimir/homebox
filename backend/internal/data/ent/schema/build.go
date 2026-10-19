package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/schema/mixins"
)

// Build holds the schema definition for the Build entity. A Build is a named
// assembly of parts (e.g. a skateboard or bike build). Parts are attached via
// BuildPart rows which keep added/removed dates, so a part's mounting history
// across builds is preserved.
type Build struct {
	ent.Schema
}

func (Build) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.BaseMixin{},
		mixins.DetailsMixin{},
		GroupMixin{ref: "builds"},
	}
}

// Fields of the Build.
func (Build) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Values(
				"assembled",
				"disassembled",
				"retired",
			).
			Default("assembled"),
		field.String("sport").
			MaxLen(64).
			Optional(),
	}
}

// Edges of the Build.
func (Build) Edges() []ent.Edge {
	owned := func(s string, t any) ent.Edge {
		return edge.To(s, t).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			})
	}

	return []ent.Edge{
		owned("parts", BuildPart.Type),
		owned("usage_sessions", UsageSession.Type),
	}
}
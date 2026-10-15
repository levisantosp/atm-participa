package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Issue holds the schema definition for the Issue entity.
type Issue struct {
	ent.Schema
}

// Fields of the Issue.
func (Issue) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("title").MaxLen(72),
		field.String("description").MaxLen(65_000),
		field.Enum("status").
			Values("open", "closed", "in_review").
			Default("open"),
		field.Int64("upvotes").Default(0),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Issue.
func (Issue) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("issues").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("issue_upvotes", Upvote.Type).Ref("issue"),
		edge.From("issue_files", File.Type).Ref("issue"),
	}
}

func (Issue) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("user"),
	}
}

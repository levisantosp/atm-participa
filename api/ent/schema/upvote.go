package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Upvote holds the schema definition for the Upvote entity.
type Upvote struct {
	ent.Schema
}

// Fields of the Upvote.
func (Upvote) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("user_id"),
		field.Int64("issue_id"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Upvote.
func (Upvote) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("issue", Issue.Type).
			Field("issue_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Upvote) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "issue_id").Unique(),
		index.Fields("issue_id"),
	}
}

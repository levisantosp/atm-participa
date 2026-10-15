package schema

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				byte := make([]byte, 32)
				if _, err := rand.Read(byte); err != nil {
					panic(err)
				}
				return hex.EncodeToString(byte)
			}).
			Immutable(),
		field.Int64("issue_id"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("issue", Issue.Type).
			Field("issue_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (File) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("issue_id"),
	}
}

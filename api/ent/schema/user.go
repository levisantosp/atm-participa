package schema

import (
	"context"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/ent/generated/hook"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("username").Unique(),
		field.String("display_name"),
		field.Bool("is_admin").Default(false),
		field.String("email").Unique(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("accounts", Account.Type),
		edge.To("issues", Issue.Type),
	}
}

func (User) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(func(m ent.Mutator) ent.Mutator {
			return hook.UserFunc(
				func(ctx context.Context, um *generated.UserMutation) (generated.Value, error) {
					if email, ok := um.Email(); ok {
						um.SetEmail(strings.ToLower(strings.TrimSpace(email)))
					}

					if username, ok := um.Username(); ok {
						um.SetUsername(
							strings.ToLower(strings.TrimSpace(username)),
						)
					}

					return m.Mutate(ctx, um)
				},
			)
		}, ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne),
	}
}

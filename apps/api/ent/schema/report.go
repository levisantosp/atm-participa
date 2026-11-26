package schema

import (
	"context"
	"errors"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/hook"
)

// Report holds the schema definition for the Report entity.
type Report struct {
	ent.Schema
}

// Fields of the Report.
func (Report) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("reporter_id").Optional().Nillable().Immutable(),
		field.Int64("target_user_id").Optional().Nillable().Immutable(),
		field.Int64("target_issue_id").Optional().Nillable().Immutable(),
		field.Enum("entity").Values("user", "issue").Immutable(),
		field.Int64("entity_id").Immutable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Report.
func (Report) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("reporter", User.Type).
			Field("reporter_id").Unique().Immutable().
			Annotations(entsql.OnDelete(entsql.SetNull)),
		edge.To("target_user", User.Type).
			Field("target_user_id").Unique().Immutable().
			Annotations(entsql.OnDelete(entsql.SetNull)),
		edge.To("target_issue", Issue.Type).
			Field("target_issue_id").Unique().Immutable().
			Annotations(entsql.OnDelete(entsql.SetNull)),
	}
}

func (Report) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Checks(map[string]string{
			"report_target_matches_entity": `(
				entity = 'user' AND target_issue_id IS NULL
				AND (target_user_id IS NULL OR target_user_id = entity_id)
			) OR (
				entity = 'issue' AND target_user_id IS NULL
				AND (target_issue_id IS NULL OR target_issue_id = entity_id)
			)`,
		}),
	}
}

// Require a live target at creation. The database can clear the edge later.
func (Report) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(
				func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
					_, userSet := m.Field("target_user_id")
					_, issueSet := m.Field("target_issue_id")
					if userSet == issueSet {
						return nil, errors.New(
							"report requires exactly one target",
						)
					}
					return next.Mutate(ctx, m)
				},
			)
		}, ent.OpCreate),
	}
}

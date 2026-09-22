package dtos

import (
	"time"

	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
)

type Issue struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      issue.Status `json:"status"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

func IssueFrom(issue *generated.Issue) Issue {
	return Issue{
		ID:          issue.ID,
		Title:       issue.Title,
		Description: issue.Description,
		Status:      issue.Status,
		CreatedAt:   issue.CreatedAt,
		UpdatedAt:   issue.UpdatedAt,
	}
}

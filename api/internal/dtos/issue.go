package dtos

import (
	"strconv"
	"time"

	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
)

type Issue struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      issue.Status `json:"status"`
	Images      []*string    `json:"images,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

func IssueFrom(issue *generated.Issue) Issue {
	images := make([]*string, 0, len(issue.Edges.IssueFiles))
	for _, file := range issue.Edges.IssueFiles {
		images = append(images, &file.ID)
	}

	return Issue{
		ID:          strconv.FormatInt(issue.ID, 10),
		Title:       issue.Title,
		Description: issue.Description,
		Status:      issue.Status,
		Images:      images,
		CreatedAt:   issue.CreatedAt,
		UpdatedAt:   issue.UpdatedAt,
	}
}

func IssuesFrom(issues []*generated.Issue) []Issue {
	items := make([]Issue, 0, len(issues))
	for _, issue := range issues {
		items = append(items, IssueFrom(issue))
	}

	return items
}

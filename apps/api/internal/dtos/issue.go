package dtos

import (
	"strconv"
	"time"

	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/issue"
)

type Issue struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	CEP          string       `json:"cep"`
	Address      string       `json:"address"`
	Neighborhood string       `json:"neighborhood"`
	Status       issue.Status `json:"status"`
	Upvotes      int64        `json:"upvotes"`
	Images       []*string    `json:"images,omitempty"`
	Author       *User        `json:"author,omitempty"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

func IssueFrom(issue *generated.Issue) Issue {
	images := make([]*string, 0, len(issue.Edges.IssueFiles))
	for _, file := range issue.Edges.IssueFiles {
		images = append(images, &file.ID)
	}

	var author *User
	if issue.Edges.User != nil {
		user := UserFrom(issue.Edges.User)
		author = &user
	}

	return Issue{
		ID:           strconv.FormatInt(issue.ID, 10),
		Title:        issue.Title,
		Description:  issue.Description,
		CEP:          issue.Cep,
		Address:      issue.Address,
		Neighborhood: issue.Neighborhood,
		Status:       issue.Status,
		Upvotes:      issue.Upvotes,
		Images:       images,
		Author:       author,
		CreatedAt:    issue.CreatedAt,
		UpdatedAt:    issue.UpdatedAt,
	}
}

func IssuesFrom(issues []*generated.Issue) []Issue {
	items := make([]Issue, 0, len(issues))
	for _, issue := range issues {
		items = append(items, IssueFrom(issue))
	}

	return items
}

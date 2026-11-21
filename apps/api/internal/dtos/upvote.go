package dtos

import (
	"strconv"
	"time"

	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
)

type Upvote struct {
	ID        string    `json:"id"`
	IssueID   string    `json:"issueId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func UpvoteFrom(upvote *generated.Upvote) Upvote {
	return Upvote{
		ID:        strconv.FormatInt(upvote.ID, 10),
		IssueID:   strconv.FormatInt(upvote.IssueID, 10),
		CreatedAt: upvote.CreatedAt,
		UpdatedAt: upvote.UpdatedAt,
	}
}

func UpvotesFrom(upvotes []*generated.Upvote) []Upvote {
	items := make([]Upvote, 0, len(upvotes))
	for _, upvote := range upvotes {
		items = append(items, UpvoteFrom(upvote))
	}
	return items
}

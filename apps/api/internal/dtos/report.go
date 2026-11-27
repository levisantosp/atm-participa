package dtos

import (
	"strconv"
	"time"

	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/report"
)

type Report struct {
	ID        string        `json:"id"`
	UserID    *string       `json:"userId"`
	Author    *User         `json:"author,omitempty"`
	Entity    report.Entity `json:"entity"`
	Reason    string        `json:"reason"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

func ReportFrom(r *generated.Report) Report {
	item := Report{
		ID:        strconv.FormatInt(r.ID, 10),
		Entity:    r.Entity,
		Reason:    r.Reason,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	if r.UserID != nil {
		id := strconv.FormatInt(*r.UserID, 10)
		item.UserID = &id
	}
	if r.Edges.User != nil {
		user := UserFrom(r.Edges.User)
		item.Author = &user
	}
	return item
}

func ReportsFrom(reports []*generated.Report) []Report {
	items := make([]Report, 0, len(reports))
	for _, r := range reports {
		items = append(items, ReportFrom(r))
	}
	return items
}

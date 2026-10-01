package routes

import (
	"context"
	"fmt"

	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/r2"
)

func IssueResponsesFrom(
	ctx context.Context,
	issues []*generated.Issue,
	client *r2.Client,
) ([]dtos.Issue, error) {
	responses := make([]dtos.Issue, 0, len(issues))

	for _, item := range issues {
		response := dtos.IssueFrom(item)
		if len(item.Edges.IssueFiles) > 0 {
			file := item.Edges.IssueFiles[0]
			imageURL, err := client.PresignGetObject(
				ctx,
				fmt.Sprintf("issues/%d/%s", item.ID, file.ID),
			)
			if err != nil {
				return nil, err
			}
			response.ImageURL = imageURL
		}

		responses = append(responses, response)
	}

	return responses, nil
}

package issues

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/dtos"
	"github.com/levisantosp/atm-participa/api/internal/middlewares"
	r2client "github.com/levisantosp/atm-participa/api/internal/r2"
	"github.com/levisantosp/atm-participa/api/internal/routes"
	"github.com/levisantosp/atm-participa/infra/r2"
)

type CreateIssueOutput struct {
	Body dtos.Issue
}

func CreateIssue(
	ctx context.Context,
	input *struct {
		RawBody huma.MultipartFormFiles[struct {
			Title       string        `form:"title" maxLength:"72" minLength:"3"`
			Description string        `form:"description" maxLength:"65000" minLength:"10"`
			File        huma.FormFile `form:"file" contentType:"image/png,image/jpeg,image/webp,image/gif" required:"false"`
		}]
	},
) (*CreateIssueOutput, error) {
	userCtx := middlewares.MustGetUserFromContext(ctx)

	body := input.RawBody.Data()

	if body.File.IsSet && body.File.Size > 15_000_000 {
		return nil, huma.Error422UnprocessableEntity(
			"O tamanho do arquivo deve ter no máximo 15 MB.",
		)
	}

	var client *r2.Client

	issue, err := db.WithTx(
		ctx,
		func(tx *generated.Tx) (*generated.Issue, error) {
			issue, err := tx.Issue.Create().
				SetTitle(body.Title).
				SetDescription(body.Description).
				SetUserID(userCtx.ID).
				Save(ctx)
			if err != nil {
				return nil, err
			}

			if body.File.IsSet {
				file, err := tx.File.Create().SetIssueID(issue.ID).Save(ctx)
				if err != nil {
					return nil, err
				}
				issue.Edges.IssueFiles = append(issue.Edges.IssueFiles, file)

				client = r2client.Client

				_, err = client.S3.PutObject(ctx, &s3.PutObjectInput{
					Bucket: aws.String(client.Bucket),
					Key: aws.String(
						fmt.Sprintf("issues/%d/%s", issue.ID, file.ID),
					),
					Body:        body.File.File,
					ContentType: aws.String(body.File.ContentType),
				})
				if err != nil {
					return nil, err
				}
			}

			return issue, nil
		},
	)
	if err != nil {
		return nil, err
	}

	response := dtos.IssueFrom(issue)
	if body.File.IsSet {
		responses, err := routes.IssueResponsesFrom(
			ctx,
			[]*generated.Issue{issue},
			client,
		)
		if err != nil {
			return nil, err
		}

		response = responses[0]
	}

	return &CreateIssueOutput{
		Body: response,
	}, nil
}

package admin

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/r2"
	"github.com/levisantosp/atm-participa/api/utils"
)

func DeleteIssue(
	ctx context.Context,
	input *struct {
		ID int64 `path:"id"`
	},
) (*struct{}, error) {
	_, err := db.WithTx(ctx, func(tx *generated.Tx) (*struct{}, error) {
		issue, err := tx.Issue.Query().
			Where(issue.IDEQ(input.ID)).
			WithIssueFiles().
			Only(ctx)
		if err != nil {
			return nil, err
		}

		client, err := r2.New(ctx)
		if err != nil {
			return nil, err
		}

		objects := make(
			[]types.ObjectIdentifier,
			0,
			len(issue.Edges.IssueFiles),
		)

		for _, file := range issue.Edges.IssueFiles {
			objects = append(objects, types.ObjectIdentifier{
				Key: aws.String(fmt.Sprintf("issues/%d/%s", issue.ID, file.ID)),
			})
		}

		if len(objects) > 0 {
			_, err := client.S3.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: aws.String(client.Bucket),
				Delete: &types.Delete{
					Objects: objects,
				},
			})
			if err != nil {
				return nil, err
			}
		}

		if err := tx.Issue.DeleteOneID(issue.ID).Exec(ctx); err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Demanda não encontrada")
		}

		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	return nil, nil
}

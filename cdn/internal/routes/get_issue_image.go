package routes

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/cdn/internal/r2"
)

func GetIssueImage(
	ctx context.Context,
	input *struct {
		ID  string `path:"id"`
		Key string `path:"key"`
	},
) (*huma.StreamResponse, error) {
	file, err := r2.Client.S3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r2.Client.Bucket),
		Key:    aws.String(fmt.Sprintf("issues/%s/%s", input.ID, input.Key)),
	})
	if err != nil {
		var e smithy.APIError
		if errors.As(err, &e) && e.ErrorCode() == "NoSuchKey" {
			return nil, huma.Error404NotFound("Imagem não encontrada")
		}

		return nil, err
	}

	contentType := "application/octet-stream"
	if file.ContentType != nil {
		contentType = *file.ContentType
	}

	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {
			ctx.SetHeader("Content-Type", contentType)
			ctx.SetHeader("Cache-Control", "public, max-age=3600")

			defer func() {
				_ = file.Body.Close()
			}()

			_, _ = io.Copy(ctx.BodyWriter(), file.Body)
		},
	}, nil
}

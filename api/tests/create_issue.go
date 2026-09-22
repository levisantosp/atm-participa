package tests

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/r2"
)

func CreateIssue(t *testing.T, user *generated.User) (*generated.Issue, error) {
	return db.Client.Issue.Create().
		SetTitle(strings.Repeat("a", 20)).
		SetDescription(strings.Repeat("a", 20)).
		SetUserID(user.ID).
		Save(t.Context())
}

func CreateIssueWithImage(
	t *testing.T,
	user *generated.User,
) (*generated.Issue, error) {
	return db.WithTx(
		t.Context(),
		func(tx *generated.Tx) (*generated.Issue, error) {
			issue, err := tx.Issue.Create().
				SetTitle(strings.Repeat("a", 20)).
				SetDescription(strings.Repeat("a", 20)).
				SetUserID(user.ID).
				Save(t.Context())
			if err != nil {
				return nil, err
			}

			file, err := tx.File.Create().
				SetIssueID(issue.ID).
				Save(t.Context())
			if err != nil {
				return nil, err
			}

			client, err := r2.New(t.Context())
			if err != nil {
				return nil, err
			}

			_, err = client.S3.PutObject(t.Context(), &s3.PutObjectInput{
				Bucket: aws.String(client.Bucket),
				Key: aws.String(
					fmt.Sprintf("issues/%d/%s", issue.ID, file.ID),
				),
				Body: bytes.NewReader(
					[]byte{
						0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
						0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
						0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
						0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
						0x89, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44,
						0x41, 0x54, 0x78, 0x9C, 0x63, 0xF8, 0xCF, 0xC0,
						0xF0, 0x1F, 0x00, 0x05, 0x00, 0x01, 0xFF, 0x89,
						0x99, 0x3D, 0x1D, 0x00, 0x00, 0x00, 0x00, 0x49,
						0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
					},
				),
				ContentType: aws.String("image/png"),
			})
			if err != nil {
				return nil, err
			}

			return issue, nil
		},
	)
}

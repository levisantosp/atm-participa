package issues

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/r2"
	"github.com/levisantosp/atm-participa/api/tests"

	_ "github.com/levisantosp/atm-participa/api/ent/generated/runtime"
)

func createMultipartBody(
	t *testing.T,
	title string,
	description string,
	fileName string,
	fileContent []byte,
) (*bytes.Buffer, string) {
	t.Helper()

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("title", title); err != nil {
		t.Fatal(err)
	}

	if err := writer.WriteField("description", description); err != nil {
		t.Fatal(err)
	}

	header := make(textproto.MIMEHeader)
	header.Set(
		"Content-Disposition",
		fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName),
	)
	header.Set("Content-Type", "image/png")

	part, err := writer.CreatePart(header)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := part.Write(fileContent); err != nil {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	return &body, writer.FormDataContentType()
}

func createBody(
	t *testing.T,
	title string,
	description string,
) (*bytes.Buffer, string) {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("title", title); err != nil {
		t.Fatal(err)
	}

	if err := writer.WriteField("description", description); err != nil {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	return &body, writer.FormDataContentType()
}

func clearBucket(t *testing.T, ctx context.Context, client r2.Client) error {
	t.Helper()

	var token *string

	for {
		list, err := client.S3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            aws.String(client.Bucket),
			ContinuationToken: token,
		})
		if err != nil {
			return err
		}

		if len(list.Contents) > 0 {
			objects := make([]types.ObjectIdentifier, 0, len(list.Contents))

			for _, object := range list.Contents {
				objects = append(objects, types.ObjectIdentifier{
					Key: object.Key,
				})
			}

			_, err = client.S3.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: aws.String(client.Bucket),
				Delete: &types.Delete{
					Objects: objects,
				},
			})
			if err != nil {
				return err
			}
		}

		if !aws.ToBool(list.IsTruncated) {
			return nil
		}

		token = list.ContinuationToken
	}
}

func TestCreateIssue(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)

	Routes(api)

	user := tests.CreateUser(t)
	session := tests.CreateSession(user, t)

	t.Run("should create an issue without image", func(t *testing.T) {
		body, contentType := createBody(
			t,
			strings.Repeat("a", 3),
			strings.Repeat("a", 10),
		)

		res := api.Post(
			"/issues",
			tests.GetCookie(session.ID),
			body,
			tests.GetContentType(contentType),
		)

		if res.Code != http.StatusCreated {
			t.Fatalf("expected status %d got %d", http.StatusCreated, res.Code)
		}

		var issue dtos.Issue
		if err := json.NewDecoder(res.Body).Decode(&issue); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should create an issue with image", func(t *testing.T) {
		client, err := r2.New(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		if err = clearBucket(t, t.Context(), *client); err != nil {
			t.Fatal(err)
		}

		body, contentType := createMultipartBody(
			t,
			strings.Repeat("a", 3),
			strings.Repeat("a", 10),
			"image.png",
			[]byte{
				0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
				0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
				0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
				0x89, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x44, 0x41,
				0x54, 0x78, 0x9C, 0x63, 0xF8, 0xCF, 0xC0, 0xF0,
				0x1F, 0x00, 0x05, 0x00, 0x01, 0xFF, 0x89, 0x99,
				0x3D, 0x1D, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45,
				0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
			},
		)

		res := api.Post(
			"/issues",
			tests.GetCookie(session.ID),
			tests.GetContentType(contentType),
			body,
		)

		if res.Code != http.StatusCreated {
			t.Fatalf("expected status %d got %d", http.StatusCreated, res.Code)
		}

		var issue dtos.Issue
		if err := json.NewDecoder(res.Body).Decode(&issue); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		body, contentType := createBody(
			t,
			strings.Repeat("a", 3),
			strings.Repeat("a", 10),
		)

		res := api.Post(
			"/issues",
			body,
			tests.GetContentType(contentType),
		)

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should reject missing title", func(t *testing.T) {
		body, contentType := createBody(
			t,
			"",
			strings.Repeat("a", 10),
		)

		res := api.Post(
			"/issues",
			tests.GetCookie(session.ID),
			body,
			tests.GetContentType(contentType),
		)

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject missing description", func(t *testing.T) {
		body, contentType := createBody(
			t,
			strings.Repeat("a", 3),
			"",
		)

		res := api.Post(
			"/issues",
			tests.GetCookie(session.ID),
			body,
			tests.GetContentType(contentType),
		)

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject title shorter than 3 characters", func(t *testing.T) {
		body, contentType := createBody(
			t,
			strings.Repeat("a", 2),
			strings.Repeat("a", 10),
		)

		res := api.Post(
			"/issues",
			tests.GetCookie(session.ID),
			body,
			tests.GetContentType(contentType),
		)

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject title longer than 72 characters", func(t *testing.T) {
		body, contentType := createBody(
			t,
			strings.Repeat("a", 73),
			strings.Repeat("a", 10),
		)

		res := api.Post(
			"/issues",
			tests.GetCookie(session.ID),
			body,
			tests.GetContentType(contentType),
		)

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run(
		"should reject description shorter than 10 characters",
		func(t *testing.T) {
			body, contentType := createBody(
				t,
				strings.Repeat("a", 3),
				strings.Repeat("a", 9),
			)

			res := api.Post(
				"/issues",
				tests.GetCookie(session.ID),
				body,
				tests.GetContentType(contentType),
			)

			if res.Code != http.StatusUnprocessableEntity {
				t.Fatalf(
					"expected status %d got %d",
					http.StatusUnprocessableEntity,
					res.Code,
				)
			}
		},
	)

	t.Run(
		"should reject description longer than 65000 characters",
		func(t *testing.T) {
			body, contentType := createBody(
				t,
				strings.Repeat("a", 3),
				strings.Repeat("a", 65_001),
			)

			res := api.Post(
				"/issues",
				tests.GetCookie(session.ID),
				body,
				tests.GetContentType(contentType),
			)

			if res.Code != http.StatusUnprocessableEntity {
				t.Fatalf(
					"expected status %d got %d",
					http.StatusUnprocessableEntity,
					res.Code,
				)
			}
		},
	)
}

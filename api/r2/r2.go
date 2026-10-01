package r2

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/levisantosp/atm-participa/api/utils"
)

type Client struct {
	S3     *s3.Client
	Bucket string
}

func New(ctx context.Context) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				utils.Env.CloudflareR2AccessKeyID,
				utils.Env.CloudflareR2SecretKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(utils.Env.CloudflareR2URL)
	})

	return &Client{
		S3:     client,
		Bucket: utils.Env.CloudflareR2Bucket,
	}, nil
}

func (c *Client) PresignGetObject(
	ctx context.Context,
	key string,
) (string, error) {
	presigned, err := s3.NewPresignClient(c.S3).PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(c.Bucket),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(15*time.Minute),
	)
	if err != nil {
		return "", err
	}

	return presigned.URL, nil
}

func (c *Client) IssueFileObjectKeys(
	ctx context.Context,
	issueID int64,
	fileID string,
) ([]string, error) {
	prefix := IssueFileObjectPrefix(issueID, fileID)
	paginator := s3.NewListObjectsV2Paginator(c.S3, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.Bucket),
		Prefix: aws.String(prefix),
	})

	keys := make([]string, 0, 1)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, object := range page.Contents {
			key := aws.ToString(object.Key)
			if key == prefix || strings.HasPrefix(key, prefix+".") {
				keys = append(keys, key)
			}
		}
	}

	sort.Strings(keys)

	return keys, nil
}

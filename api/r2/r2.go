package r2

import (
	"context"

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

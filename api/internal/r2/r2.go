package r2

import (
	"context"

	"github.com/levisantosp/atm-participa/api/internal/env"
	"github.com/levisantosp/atm-participa/infra/r2"
)

var Client *r2.Client

func Connect(ctx context.Context) error {
	client, err := r2.New(ctx, r2.Config{
		AccessKeyID: env.Data.CloudflareR2AccessKeyID,
		SecretKey:   env.Data.CloudflareR2SecretKey,
		URL:         env.Data.CloudflareR2URL,
		Bucket:      env.Data.CloudflareR2Bucket,
	})
	if err != nil {
		return err
	}

	Client = client
	return nil
}

package r2

import (
	"context"

	"github.com/levisantosp/atm-participa/api/internal/env"
	"github.com/levisantosp/atm-participa/infra/r2"
)

var Client *r2.Client

func Connect(ctx context.Context) error {
	client, err := r2.New(ctx, r2.Config{
		AccessKeyID: env.Env.CloudflareR2AccessKeyID,
		SecretKey:   env.Env.CloudflareR2SecretKey,
		URL:         env.Env.CloudflareR2URL,
		Bucket:      env.Env.CloudflareR2Bucket,
	})
	if err != nil {
		return err
	}

	Client = client
	return nil
}

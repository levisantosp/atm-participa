package r2

import (
	"context"

	"github.com/levisantosp/atm-participa/api/internal/utils"
	"github.com/levisantosp/atm-participa/infra/r2"
)

var Client *r2.Client

func Connect(ctx context.Context) error {
	client, err := r2.New(ctx, r2.Config{
		AccessKeyID: utils.Env.CloudflareR2AccessKeyID,
		SecretKey:   utils.Env.CloudflareR2SecretKey,
		URL:         utils.Env.CloudflareR2URL,
		Bucket:      utils.Env.CloudflareR2Bucket,
	})
	if err != nil {
		return err
	}

	Client = client
	return nil
}

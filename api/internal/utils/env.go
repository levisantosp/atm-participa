package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zenv"
	"github.com/joho/godotenv"
)

type TEnv struct {
	DatabaseURL             string   `zog:"DATABASE_URL"`
	RedisAddr               string   `zog:"REDIS_ADDR"`
	RedisPassword           string   `zog:"REDIS_PASSWORD"`
	TrustedOrigins          []string `zog:"TRUSTED_ORIGINS"`
	CloudflareR2AccessKeyID string   `zog:"CLOUDFLARE_R2_ACCESS_KEY_ID"`
	CloudflareR2SecretKey   string   `zog:"CLOUDFLARE_R2_SECRET_ACCESS_KEY"`
	CloudflareR2URL         string   `zog:"CLOUDFLARE_R2_URL"`
	CloudflareR2Bucket      string   `zog:"CLOUDFLARE_R2_BUCKET"`
}

var Env TEnv

func LoadEnv(envFile string) {
	_ = godotenv.Overload(envFile)

	schema := zog.Struct(zog.Shape{
		"DatabaseURL":   zog.String().URL().Required(),
		"RedisAddr":     zog.String().Required(),
		"RedisPassword": zog.String().Required(),
		"TrustedOrigins": zog.Preprocess(
			func(data any, ctx zog.Ctx) (any, error) {
				value, ok := data.(string)
				if !ok {
					return nil, fmt.Errorf("expected string, got %T", data)
				}

				return strings.Split(value, ","), nil
			},
			zog.Slice(zog.String().URL().Required()),
		),
		"CloudflareR2AccessKeyID": zog.String().Required(),
		"CloudflareR2SecretKey":   zog.String().Required(),
		"CloudflareR2URL":         zog.String().URL().Required(),
		"CloudflareR2Bucket":      zog.String().Required(),
	})

	err := schema.Parse(zenv.NewDataProvider(), &Env)
	if err != nil {
		log.Fatal(err)
	}
}

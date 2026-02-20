package config

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetR2Client() *s3.Client {
	accountId := os.Getenv("ACCOUNT_ID")
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("SECRET_ACCESS_KEY")

	if accountId == "" || accessKeyId == "" || accessKeySecret == "" {
		log.Panic().Msg("u2 credentials shouldn't be empty")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId)

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"), // Required by SDK but not used by R2
	)
	if err != nil {
		log.Panic().Err(err).Msg("unable to load AWS SDK config")
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return client
}

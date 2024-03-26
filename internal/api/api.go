//go:generate go run ec2/ec2_gen.go
//go:generate go run iam/iam_gen.go
//go:generate go run s3/s3_gen.go
//go:generate gofmt -w ec2/
//go:generate gofmt -w iam/
//go:generate gofmt -w s3/

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const DefaultRegion = "us-east-1"

var DefaultTargetRegions = []string{
	"ap-south-1",
	"eu-north-1",
	"eu-west-3",
	"eu-west-2",
	"eu-west-1",
	"ap-northeast-3",
	"ap-northeast-2",
	"ap-northeast-1",
	"ca-central-1",
	"sa-east-1",
	"ap-southeast-1",
	"ap-southeast-2",
	"eu-central-1",
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
}

func LoadConfig(ctx context.Context, region string, profile string) (*aws.Config, error) {
	var cfg aws.Config
	var err error
	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot load aws config: %w", err)
	}
	if region != "" {
		cfg.Region = region
	}
	if cfg.Region == "" {
		cfg.Region = DefaultRegion
	}
	cfg.RetryMode = aws.RetryModeStandard
	cfg.RetryMaxAttempts = 10
	return &cfg, nil
}

func DecodePolicyDocument(document string, unescape bool) (string, error) {
	if unescape {
		decoded, err := url.QueryUnescape(document)
		if err != nil {
			return "", err
		}
		document = decoded
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(document), "", "  "); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func Contains(s string, targets []string) bool {
	for _, target := range targets {
		if strings.Contains(s, target) {
			return true
		}
	}
	return false
}

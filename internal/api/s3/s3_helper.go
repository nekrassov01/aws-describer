package s3

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nekrassov01/aws-describer/internal/api"
)

func getS3BucketLocation(ctx context.Context, client *s3.Client, name *string) (string, bool, error) {
	o, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{Bucket: name})
	if err != nil {
		if strings.Contains(err.Error(), "AccessDenied") {
			return "", false, nil
		}
		return "", false, err
	}
	if o.LocationConstraint == "" {
		return "us-east-1", true, nil
	}
	return string(o.LocationConstraint), true, nil
}

func getS3BucketPolicyDocument(ctx context.Context, client *s3.Client, name *string, region string) (string, bool, error) {
	o, err := client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
		Bucket: name,
	}, func(opt *s3.Options) {
		opt.Region = region
	})
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchBucketPolicy") {
			return "", false, nil
		}
		return "", false, err
	}
	doc, err := api.DecodePolicyDocument(aws.ToString(o.Policy), false)
	if err != nil {
		return "", false, err
	}
	return doc, true, nil
}

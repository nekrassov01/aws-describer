package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var _ IS3Client = (*S3Client)(nil)

type IS3Client interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)

	GetBucketLocation(ctx context.Context, name *string) (string, bool, error)
	GetBucketPolicyDocument(ctx context.Context, name *string, region string) (string, bool, error)
}

type S3Client struct {
	*s3.Client
}

func NewS3Client(cfg *aws.Config) *S3Client {
	return &S3Client{Client: s3.NewFromConfig(*cfg)}
}

func (client *S3Client) GetBucketLocation(ctx context.Context, name *string) (string, bool, error) {
	return getS3BucketLocation(ctx, client.Client, name)
}

func (client *S3Client) GetBucketPolicyDocument(ctx context.Context, name *string, region string) (string, bool, error) {
	return getS3BucketPolicyDocument(ctx, client.Client, name, region)
}

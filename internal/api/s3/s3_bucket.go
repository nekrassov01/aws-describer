package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/nekrassov01/aws-describer/internal/api"
)

type BucketInfo struct {
	BucketName     string
	IsAccesible    bool
	Location       string
	PolicyDocument string
}

func GetBucketInfo(ctx context.Context, client IS3Client, ich chan<- BucketInfo, bucket types.Bucket, document bool, filters []string) error {
	region, isAccesible, err := client.GetBucketLocation(ctx, bucket.Name)
	if err != nil {
		return err
	}
	if !isAccesible {
		if len(filters) > 0 {
			return nil
		}
		appendToBucketInfo(ich, bucket, isAccesible, "", "")
		return nil
	}
	var doc string
	var hasPolicy bool
	if document {
		doc, hasPolicy, err = client.GetBucketPolicyDocument(ctx, bucket.Name, region)
		if err != nil {
			return err
		}
	}
	if !hasPolicy {
		if len(filters) > 0 {
			return nil
		}
		appendToBucketInfo(ich, bucket, isAccesible, region, "")
		return nil
	}
	if len(filters) == 0 || api.Contains(doc, filters) {
		appendToBucketInfo(ich, bucket, isAccesible, region, doc)
	}
	return nil
}

func appendToBucketInfo(ich chan<- BucketInfo, bucket types.Bucket, isAcccesible bool, region, doc string) {
	ich <- BucketInfo{
		BucketName:     aws.ToString(bucket.Name),
		IsAccesible:    isAcccesible,
		Location:       region,
		PolicyDocument: doc,
	}
}

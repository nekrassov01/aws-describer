package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/nekrassov01/aws-describer/internal/api"
)

type PolicyInfo struct {
	PolicyName                    string
	PolicyId                      string
	PolicyArn                     string
	Path                          string
	DefaultVersionId              string
	IsAttachable                  bool
	AttachmentCount               int32
	PermissionsBoundaryUsageCount int32
	PolicyDocument                string
}

func GetPolicyInfo(ctx context.Context, client IIamClient, ich chan<- PolicyInfo, policy types.Policy, document bool, filters []string) error {
	var doc string
	var err error
	if document {
		doc, err = client.GetPolicyDocument(ctx, policy.Arn, policy.DefaultVersionId)
		if err != nil {
			return err
		}
	}
	if len(filters) == 0 || api.Contains(doc, filters) {
		ich <- PolicyInfo{
			PolicyName:                    aws.ToString(policy.PolicyName),
			PolicyId:                      aws.ToString(policy.PolicyId),
			PolicyArn:                     aws.ToString(policy.Arn),
			Path:                          aws.ToString(policy.Path),
			DefaultVersionId:              aws.ToString(policy.DefaultVersionId),
			IsAttachable:                  policy.IsAttachable,
			AttachmentCount:               aws.ToInt32(policy.AttachmentCount),
			PermissionsBoundaryUsageCount: aws.ToInt32(policy.PermissionsBoundaryUsageCount),
			PolicyDocument:                doc,
		}
	}
	return nil
}

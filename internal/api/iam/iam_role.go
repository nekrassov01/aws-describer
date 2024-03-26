package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/nekrassov01/aws-describer/internal/api"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

type RoleInfo struct {
	RoleName string
	RoleId   string
	Path     string
	RoleArn  string
}

func GetRoleInfo(ich chan<- RoleInfo, role types.Role) {
	ich <- RoleInfo{
		RoleName: aws.ToString(role.RoleName),
		RoleId:   aws.ToString(role.RoleId),
		Path:     aws.ToString(role.Path),
		RoleArn:  aws.ToString(role.Arn),
	}
}

type RolePolicyInfo struct {
	RoleName       string
	RoleId         string
	Path           string
	PolicyType     string
	PolicyName     string
	PolicyDocument string
}

func GetRolePolicyInfo(ctx context.Context, l *rate.Limiter, client IIamClient, ich chan<- RolePolicyInfo, role types.Role, document bool, filters []string, pols map[string]types.Policy) error {
	var hasAttachedPolicy, hasInlinePolicy bool
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		hasAttachedPolicy, err = appendToRolePolicyInfoForAttachedPolicy(ctx, client, ich, role, document, filters, pols)
		return err
	})
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		hasInlinePolicy, err = appendToRolePolicyInfoForInlinePolicy(ctx, client, ich, role, document, filters)
		return err
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	if !hasAttachedPolicy && !hasInlinePolicy {
		if len(filters) > 0 {
			return nil
		}
		appendToRolePolicyInfo(ich, role, "", "", "")
	}
	return nil
}

func appendToRolePolicyInfoForAttachedPolicy(ctx context.Context, client IIamClient, ich chan<- RolePolicyInfo, role types.Role, document bool, filters []string, pols map[string]types.Policy) (bool, error) {
	found := false
	apols, err := client.GetAttachedRolePolicies(ctx, role.RoleName)
	if err != nil {
		return false, err
	}
	for _, apol := range apols {
		found = true
		if !document {
			appendToRolePolicyInfo(ich, role, policyTypeAttached.String(), aws.ToString(apol.PolicyName), "")
			continue
		}
		doc, err := client.GetCustomerPolicyDocument(ctx, apol.PolicyArn, pols)
		if err != nil {
			return false, err
		}
		if len(filters) == 0 || api.Contains(doc, filters) {
			appendToRolePolicyInfo(ich, role, policyTypeAttached.String(), aws.ToString(apol.PolicyName), doc)
		}
	}
	return found, nil
}

func appendToRolePolicyInfoForInlinePolicy(ctx context.Context, client IIamClient, ich chan<- RolePolicyInfo, role types.Role, document bool, filters []string) (bool, error) {
	found := false
	ipols, err := client.GetInlineRolePolicies(ctx, role.RoleName)
	if err != nil {
		return false, err
	}
	for _, ipol := range ipols {
		found = true
		if !document {
			appendToRolePolicyInfo(ich, role, policyTypeInline.String(), ipol, "")
			continue
		}
		doc, err := client.GetInlineRolePolicyDocument(ctx, role.RoleName, aws.String(ipol))
		if err != nil {
			return false, err
		}
		if len(filters) == 0 || api.Contains(doc, filters) {
			appendToRolePolicyInfo(ich, role, policyTypeInline.String(), ipol, doc)
		}
	}
	return found, nil
}

func appendToRolePolicyInfo(ich chan<- RolePolicyInfo, role types.Role, policyType, policyName, policyDocument string) {
	ich <- RolePolicyInfo{
		RoleName:       aws.ToString(role.RoleName),
		RoleId:         aws.ToString(role.RoleId),
		Path:           aws.ToString(role.Path),
		PolicyType:     policyType,
		PolicyName:     policyName,
		PolicyDocument: policyDocument,
	}
}

type RoleAssumeInfo struct {
	RoleName                 string
	RoleId                   string
	Path                     string
	AssumeRolePolicyDocument string
}

func GetRoleAssumeInfo(ich chan<- RoleAssumeInfo, role types.Role) error {
	doc, err := api.DecodePolicyDocument(aws.ToString(role.AssumeRolePolicyDocument), true)
	if err != nil {
		return err
	}
	ich <- RoleAssumeInfo{
		RoleName:                 aws.ToString(role.RoleName),
		RoleId:                   aws.ToString(role.RoleId),
		Path:                     aws.ToString(role.Path),
		AssumeRolePolicyDocument: doc,
	}
	return nil
}

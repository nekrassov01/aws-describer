package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/nekrassov01/aws-describer/internal/api"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

type GroupInfo struct {
	GroupName string
	GroupId   string
	Path      string
	GroupArn  string
}

func GetGroupInfo(ich chan<- GroupInfo, group types.Group) {
	ich <- GroupInfo{
		GroupName: aws.ToString(group.GroupName),
		GroupId:   aws.ToString(group.GroupId),
		Path:      aws.ToString(group.Path),
		GroupArn:  aws.ToString(group.Arn),
	}
}

type GroupPolicyInfo struct {
	GroupName      string
	GroupId        string
	Path           string
	PolicyType     string
	PolicyName     string
	PolicyDocument string
}

func GetGroupPolicyInfo(ctx context.Context, l *rate.Limiter, client IIamClient, ich chan<- GroupPolicyInfo, group types.Group, document bool, filters []string, pols map[string]types.Policy) error {
	var hasAttachedPolicy, hasInlinePolicy bool
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		hasAttachedPolicy, err = appendToGroupPolicyInfoForAttachedPolicy(ctx, client, ich, group, document, filters, pols)
		return err
	})
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		hasInlinePolicy, err = appendToGroupPolicyInfoForInlinePolicy(ctx, client, ich, group, document, filters)
		return err
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	if !hasAttachedPolicy && !hasInlinePolicy {
		if len(filters) > 0 {
			return nil
		}
		appendToGroupPolicyInfo(ich, group, "", "", "")
	}
	return nil
}

func appendToGroupPolicyInfoForAttachedPolicy(ctx context.Context, client IIamClient, ich chan<- GroupPolicyInfo, group types.Group, document bool, filters []string, pols map[string]types.Policy) (bool, error) {
	found := false
	apols, err := client.GetAttachedGroupPolicies(ctx, group.GroupName)
	if err != nil {
		return false, err
	}
	for _, apol := range apols {
		found = true
		if !document {
			appendToGroupPolicyInfo(ich, group, policyTypeAttached.String(), aws.ToString(apol.PolicyName), "")
			continue
		}
		doc, err := client.GetCustomerPolicyDocument(ctx, apol.PolicyArn, pols)
		if err != nil {
			return false, err
		}
		if len(filters) == 0 || api.Contains(doc, filters) {
			appendToGroupPolicyInfo(ich, group, policyTypeAttached.String(), aws.ToString(apol.PolicyName), doc)
		}
	}
	return found, nil
}

func appendToGroupPolicyInfoForInlinePolicy(ctx context.Context, client IIamClient, ich chan<- GroupPolicyInfo, group types.Group, document bool, filters []string) (bool, error) {
	found := false
	ipols, err := client.GetInlineGroupPolicies(ctx, group.GroupName)
	if err != nil {
		return false, err
	}
	for _, ipol := range ipols {
		found = true
		if !document {
			appendToGroupPolicyInfo(ich, group, policyTypeInline.String(), ipol, "")
			continue
		}
		doc, err := client.GetInlineGroupPolicyDocument(ctx, group.GroupName, aws.String(ipol))
		if err != nil {
			return false, err
		}
		if len(filters) == 0 || api.Contains(doc, filters) {
			appendToGroupPolicyInfo(ich, group, policyTypeInline.String(), ipol, doc)
		}
	}
	return found, nil
}

func appendToGroupPolicyInfo(ich chan<- GroupPolicyInfo, role types.Group, policyType, policyName, policyDocument string) {
	ich <- GroupPolicyInfo{
		GroupName:      aws.ToString(role.GroupName),
		GroupId:        aws.ToString(role.GroupId),
		Path:           aws.ToString(role.Path),
		PolicyType:     policyType,
		PolicyName:     policyName,
		PolicyDocument: policyDocument,
	}
}

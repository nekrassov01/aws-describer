package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/nekrassov01/aws-describer/internal/api"
)

func fetchIamCustomerPolicies(ctx context.Context, client *iam.Client) (map[string]types.Policy, error) {
	p := iam.NewListPoliciesPaginator(client, &iam.ListPoliciesInput{
		Scope:        types.PolicyScopeTypeLocal,
		OnlyAttached: true,
	})
	res := make(map[string]types.Policy)
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, policy := range page.Policies {
			res[aws.ToString(policy.Arn)] = policy
		}
	}
	return res, nil
}

func getIamGroupsForUser(ctx context.Context, client *iam.Client, name *string) ([]types.Group, error) {
	p := iam.NewListGroupsForUserPaginator(client, &iam.ListGroupsForUserInput{UserName: name})
	var res []types.Group
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		res = append(res, page.Groups...)
	}
	return res, nil
}

func getIamAttachedUserPolicies(ctx context.Context, client *iam.Client, name *string) ([]types.AttachedPolicy, error) {
	p := iam.NewListAttachedUserPoliciesPaginator(client, &iam.ListAttachedUserPoliciesInput{UserName: name})
	var res []types.AttachedPolicy
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		res = append(res, page.AttachedPolicies...)
	}
	return res, nil
}

func getIamInlineUserPolicies(ctx context.Context, client *iam.Client, name *string) ([]string, error) {
	p := iam.NewListUserPoliciesPaginator(client, &iam.ListUserPoliciesInput{UserName: name})
	var res []string
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		res = append(res, page.PolicyNames...)
	}
	return res, nil
}

func getIamInlineUserPolicyDocument(ctx context.Context, client *iam.Client, name *string, policyName *string) (string, error) {
	o, err := client.GetUserPolicy(ctx, &iam.GetUserPolicyInput{UserName: name, PolicyName: policyName})
	if err != nil {
		return "", err
	}
	doc, err := api.DecodePolicyDocument(aws.ToString(o.PolicyDocument), true)
	if err != nil {
		return "", err
	}
	return doc, nil
}

func getIamAttachedGroupPolicies(ctx context.Context, client *iam.Client, name *string) ([]types.AttachedPolicy, error) {
	p := iam.NewListAttachedGroupPoliciesPaginator(client, &iam.ListAttachedGroupPoliciesInput{GroupName: name})
	var res []types.AttachedPolicy
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		res = append(res, page.AttachedPolicies...)
	}
	return res, nil
}

func getIamInlineGroupPolicies(ctx context.Context, client *iam.Client, name *string) ([]string, error) {
	p := iam.NewListGroupPoliciesPaginator(client, &iam.ListGroupPoliciesInput{GroupName: name})
	var res []string
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		res = append(res, page.PolicyNames...)
	}
	return res, nil
}

func getIamInlineGroupPolicyDocument(ctx context.Context, client *iam.Client, name *string, policyName *string) (string, error) {
	o, err := client.GetGroupPolicy(ctx, &iam.GetGroupPolicyInput{GroupName: name, PolicyName: policyName})
	if err != nil {
		return "", err
	}
	doc, err := api.DecodePolicyDocument(aws.ToString(o.PolicyDocument), true)
	if err != nil {
		return "", err
	}
	return doc, nil
}

func getIamAttachedRolePolicies(ctx context.Context, client *iam.Client, name *string) ([]types.AttachedPolicy, error) {
	p := iam.NewListAttachedRolePoliciesPaginator(client, &iam.ListAttachedRolePoliciesInput{RoleName: name})
	var res []types.AttachedPolicy
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		res = append(res, page.AttachedPolicies...)
	}
	return res, nil
}

func getIamInlineRolePolicies(ctx context.Context, client *iam.Client, name *string) ([]string, error) {
	p := iam.NewListRolePoliciesPaginator(client, &iam.ListRolePoliciesInput{RoleName: name})
	var res []string
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		res = append(res, page.PolicyNames...)
	}
	return res, nil
}

func getIamInlineRolePolicyDocument(ctx context.Context, client *iam.Client, name *string, policyName *string) (string, error) {
	o, err := client.GetRolePolicy(ctx, &iam.GetRolePolicyInput{RoleName: name, PolicyName: policyName})
	if err != nil {
		return "", err
	}
	doc, err := api.DecodePolicyDocument(aws.ToString(o.PolicyDocument), true)
	if err != nil {
		return "", err
	}
	return doc, nil
}

func getIamPolicyScope(scope string) (types.PolicyScopeType, error) {
	switch scope {
	case PolicyScopeTypeLocal.String():
		return types.PolicyScopeTypeLocal, nil
	case PolicyScopeTypeAws.String():
		return types.PolicyScopeTypeAws, nil
	default:
		return "", fmt.Errorf("invalid value: %s: valid value %s: ", scope, strings.Join(PolicyScopeTypes, "|"))
	}
}

func getIamPolicyDocument(ctx context.Context, client *iam.Client, arn *string, version *string) (string, error) {
	o, err := client.GetPolicyVersion(ctx, &iam.GetPolicyVersionInput{
		PolicyArn: arn,
		VersionId: version,
	})
	if err != nil {
		return "", err
	}
	doc, err := api.DecodePolicyDocument(aws.ToString(o.PolicyVersion.Document), true)
	if err != nil {
		return "", err
	}
	return doc, nil
}

func getIamCustomerPolicyDocument(ctx context.Context, client *iam.Client, arn *string, pols map[string]types.Policy) (string, error) {
	item, ok := pols[aws.ToString(arn)]
	if !ok {
		return "SKIPPED", nil
	}
	d, err := getIamPolicyDocument(ctx, client, arn, item.DefaultVersionId)
	if err != nil {
		return "", err
	}
	return d, nil
}

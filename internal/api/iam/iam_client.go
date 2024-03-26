package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

var _ IIamClient = (*IamClient)(nil)

type IIamClient interface {
	ListUsers(ctx context.Context, params *iam.ListUsersInput, optFns ...func(*iam.Options)) (*iam.ListUsersOutput, error)
	ListGroups(ctx context.Context, params *iam.ListGroupsInput, optFns ...func(*iam.Options)) (*iam.ListGroupsOutput, error)
	ListRoles(ctx context.Context, params *iam.ListRolesInput, optFns ...func(*iam.Options)) (*iam.ListRolesOutput, error)
	ListPolicies(ctx context.Context, params *iam.ListPoliciesInput, optFns ...func(*iam.Options)) (*iam.ListPoliciesOutput, error)

	FetchCustomerPolicies(ctx context.Context) (map[string]types.Policy, error)
	GetGroupsForUser(ctx context.Context, name *string) ([]types.Group, error)
	GetAttachedUserPolicies(ctx context.Context, name *string) ([]types.AttachedPolicy, error)
	GetInlineUserPolicies(ctx context.Context, name *string) ([]string, error)
	GetInlineUserPolicyDocument(ctx context.Context, name *string, policyName *string) (string, error)
	GetAttachedGroupPolicies(ctx context.Context, name *string) ([]types.AttachedPolicy, error)
	GetInlineGroupPolicies(ctx context.Context, name *string) ([]string, error)
	GetInlineGroupPolicyDocument(ctx context.Context, name *string, policyName *string) (string, error)
	GetAttachedRolePolicies(ctx context.Context, name *string) ([]types.AttachedPolicy, error)
	GetInlineRolePolicies(ctx context.Context, name *string) ([]string, error)
	GetInlineRolePolicyDocument(ctx context.Context, name *string, policyName *string) (string, error)
	GetPolicyScope(scope string) (types.PolicyScopeType, error)
	GetPolicyDocument(ctx context.Context, arn *string, version *string) (string, error)
	GetCustomerPolicyDocument(ctx context.Context, arn *string, pols map[string]types.Policy) (string, error)
}

type IamClient struct {
	*iam.Client
}

func NewIamClient(cfg *aws.Config) *IamClient {
	return &IamClient{Client: iam.NewFromConfig(*cfg)}
}

func (client *IamClient) FetchCustomerPolicies(ctx context.Context) (map[string]types.Policy, error) {
	return fetchIamCustomerPolicies(ctx, client.Client)
}

func (client *IamClient) GetGroupsForUser(ctx context.Context, name *string) ([]types.Group, error) {
	return getIamGroupsForUser(ctx, client.Client, name)
}

func (client *IamClient) GetAttachedUserPolicies(ctx context.Context, name *string) ([]types.AttachedPolicy, error) {
	return getIamAttachedUserPolicies(ctx, client.Client, name)
}

func (client *IamClient) GetInlineUserPolicies(ctx context.Context, name *string) ([]string, error) {
	return getIamInlineUserPolicies(ctx, client.Client, name)
}

func (client *IamClient) GetInlineUserPolicyDocument(ctx context.Context, name *string, policyName *string) (string, error) {
	return getIamInlineUserPolicyDocument(ctx, client.Client, name, policyName)
}

func (client *IamClient) GetAttachedGroupPolicies(ctx context.Context, name *string) ([]types.AttachedPolicy, error) {
	return getIamAttachedGroupPolicies(ctx, client.Client, name)
}

func (client *IamClient) GetInlineGroupPolicies(ctx context.Context, name *string) ([]string, error) {
	return getIamInlineGroupPolicies(ctx, client.Client, name)
}

func (client *IamClient) GetInlineGroupPolicyDocument(ctx context.Context, name *string, policyName *string) (string, error) {
	return getIamInlineGroupPolicyDocument(ctx, client.Client, name, policyName)
}

func (client *IamClient) GetAttachedRolePolicies(ctx context.Context, name *string) ([]types.AttachedPolicy, error) {
	return getIamAttachedRolePolicies(ctx, client.Client, name)
}

func (client *IamClient) GetInlineRolePolicies(ctx context.Context, name *string) ([]string, error) {
	return getIamInlineRolePolicies(ctx, client.Client, name)
}

func (client *IamClient) GetInlineRolePolicyDocument(ctx context.Context, name *string, policyName *string) (string, error) {
	return getIamInlineRolePolicyDocument(ctx, client.Client, name, policyName)
}

func (client *IamClient) GetPolicyScope(scope string) (types.PolicyScopeType, error) {
	return getIamPolicyScope(scope)
}

func (client *IamClient) GetPolicyDocument(ctx context.Context, arn *string, version *string) (string, error) {
	return getIamPolicyDocument(ctx, client.Client, arn, version)
}

func (client *IamClient) GetCustomerPolicyDocument(ctx context.Context, arn *string, pols map[string]types.Policy) (string, error) {
	return getIamCustomerPolicyDocument(ctx, client.Client, arn, pols)
}

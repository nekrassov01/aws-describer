package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/nekrassov01/aws-describer/internal/api"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

type UserInfo struct {
	UserName string
	UserId   string
	Path     string
	UserArn  string
}

func GetUserInfo(ich chan<- UserInfo, user types.User) {
	ich <- UserInfo{
		UserName: aws.ToString(user.UserName),
		UserId:   aws.ToString(user.UserId),
		Path:     aws.ToString(user.Path),
		UserArn:  aws.ToString(user.Arn),
	}
}

type UserPolicyInfo struct {
	UserName       string
	UserId         string
	Path           string
	PolicyType     string
	PolicyName     string
	PolicyDocument string
}

func GetUserPolicyInfo(ctx context.Context, l *rate.Limiter, client IIamClient, ich chan<- UserPolicyInfo, user types.User, document bool, filters []string, pols map[string]types.Policy) error {
	var hasAttachedPolicy, hasInlinePolicy bool
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		hasAttachedPolicy, err = appendToUserPolicyInfoForAttachedPolicy(ctx, client, ich, user, document, filters, pols)
		return err
	})
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		hasInlinePolicy, err = appendToUserPolicyInfoForInlinePolicy(ctx, client, ich, user, document, filters)
		return err
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	if !hasAttachedPolicy && !hasInlinePolicy {
		if len(filters) > 0 {
			return nil
		}
		appendToUserPolicyInfo(ich, user, "", "", "")
	}
	return nil
}

func appendToUserPolicyInfoForAttachedPolicy(ctx context.Context, client IIamClient, ich chan<- UserPolicyInfo, user types.User, document bool, filters []string, pols map[string]types.Policy) (bool, error) {
	found := false
	apols, err := client.GetAttachedUserPolicies(ctx, user.UserName)
	if err != nil {
		return false, err
	}
	for _, apol := range apols {
		found = true
		if !document {
			appendToUserPolicyInfo(ich, user, policyTypeAttached.String(), aws.ToString(apol.PolicyName), "")
			continue
		}
		doc, err := client.GetCustomerPolicyDocument(ctx, apol.PolicyArn, pols)
		if err != nil {
			return false, err
		}
		if len(filters) == 0 || api.Contains(doc, filters) {
			appendToUserPolicyInfo(ich, user, policyTypeAttached.String(), aws.ToString(apol.PolicyName), doc)
		}
	}
	return found, nil
}

func appendToUserPolicyInfoForInlinePolicy(ctx context.Context, client IIamClient, ich chan<- UserPolicyInfo, user types.User, document bool, filters []string) (bool, error) {
	found := false
	ipols, err := client.GetInlineUserPolicies(ctx, user.UserName)
	if err != nil {
		return false, err
	}
	for _, ipol := range ipols {
		found = true
		if !document {
			appendToUserPolicyInfo(ich, user, policyTypeInline.String(), ipol, "")
			continue
		}
		doc, err := client.GetInlineUserPolicyDocument(ctx, user.UserName, aws.String(ipol))
		if err != nil {
			return false, err
		}
		if len(filters) == 0 || api.Contains(doc, filters) {
			appendToUserPolicyInfo(ich, user, policyTypeInline.String(), ipol, doc)
		}
	}
	return found, nil
}

func appendToUserPolicyInfo(ich chan<- UserPolicyInfo, user types.User, policyType, policyName, policyDocument string) {
	ich <- UserPolicyInfo{
		UserName:       aws.ToString(user.UserName),
		UserId:         aws.ToString(user.UserId),
		Path:           aws.ToString(user.Path),
		PolicyType:     policyType,
		PolicyName:     policyName,
		PolicyDocument: policyDocument,
	}
}

type UserGroupInfo struct {
	UserName  string
	UserId    string
	Path      string
	GroupName string
	GroupId   string
}

func GetUserGroupInfo(ctx context.Context, client IIamClient, ich chan<- UserGroupInfo, user types.User) error {
	groups, err := client.GetGroupsForUser(ctx, user.UserName)
	if err != nil {
		return err
	}
	for _, group := range groups {
		appendToUserGroupInfo(ich, user, aws.ToString(group.GroupName), aws.ToString(group.GroupId))
	}
	if len(groups) == 0 {
		appendToUserGroupInfo(ich, user, "", "")
	}
	return nil
}

func appendToUserGroupInfo(ich chan<- UserGroupInfo, user types.User, groupName, groupId string) {
	ich <- UserGroupInfo{
		UserName:  aws.ToString(user.UserName),
		UserId:    aws.ToString(user.UserId),
		Path:      aws.ToString(user.Path),
		GroupName: groupName,
		GroupId:   groupId,
	}
}

type UserAssociationInfo struct {
	UserName       string
	AttachedBy     string
	PolicyType     string
	PolicyName     string
	PolicyDocument string
}

func GetUserAssociationInfo(ctx context.Context, l *rate.Limiter, client IIamClient, ich chan<- UserAssociationInfo, user types.User, document bool, filters []string, pols map[string]types.Policy) error {
	groups, err := getUserAssociationInfoForUser(ctx, l, client, ich, user, document, filters, pols)
	if err != nil {
		return err
	}
	for _, group := range groups {
		if err := getUserAssociationInfoForGroup(ctx, l, client, ich, user, group, document, filters, pols); err != nil {
			return err
		}
	}
	return nil
}

func getUserAssociationInfoForUser(ctx context.Context, l *rate.Limiter, client IIamClient, ich chan<- UserAssociationInfo, user types.User, document bool, filters []string, pols map[string]types.Policy) ([]types.Group, error) {
	var groups []types.Group
	var hasAttachedPolicy, hasInlinePolicy bool
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		hasAttachedPolicy, err = appendToUserAssociationInfoForAttachedUserPolicy(ctx, client, ich, user, document, filters, pols)
		return err
	})
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		hasInlinePolicy, err = appendToUserAssociationInfoForInlineUserPolicy(ctx, client, ich, user, document, filters)
		return err
	})
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		groups, err = client.GetGroupsForUser(ctx, user.UserName)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	if !hasAttachedPolicy && !hasInlinePolicy {
		if len(filters) > 0 {
			return []types.Group{}, nil
		}
		appendToUserAssociationInfo(ich, aws.ToString(user.UserName), aws.ToString(user.Arn), "", "", "")
	}
	return groups, nil
}

func getUserAssociationInfoForGroup(ctx context.Context, l *rate.Limiter, client IIamClient, ich chan<- UserAssociationInfo, user types.User, group types.Group, document bool, filters []string, pols map[string]types.Policy) error {
	var hasAttachedPolicy, hasInlinePolicy bool
	geg, gctx := errgroup.WithContext(ctx)
	geg.Go(func() error {
		if err := l.Wait(gctx); err != nil {
			return err
		}
		var err error
		hasAttachedPolicy, err = appendToUserAssociationInfoForAttachedGroupPolicy(gctx, client, ich, user, group, document, filters, pols)
		return err
	})
	geg.Go(func() error {
		if err := l.Wait(gctx); err != nil {
			return err
		}
		var err error
		hasInlinePolicy, err = appendToUserAssociationInfoForInlineGroupPolicy(gctx, client, ich, user, group, document, filters)
		return err
	})
	if err := geg.Wait(); err != nil {
		return err
	}
	if !hasAttachedPolicy && !hasInlinePolicy {
		if len(filters) > 0 {
			return nil
		}
		appendToUserAssociationInfo(ich, aws.ToString(user.UserName), aws.ToString(group.Arn), "", "", "")
	}
	return nil
}

func appendToUserAssociationInfoForAttachedUserPolicy(ctx context.Context, client IIamClient, ich chan<- UserAssociationInfo, user types.User, document bool, filters []string, pols map[string]types.Policy) (bool, error) {
	found := false
	apols, err := client.GetAttachedUserPolicies(ctx, user.UserName)
	if err != nil {
		return false, err
	}
	for _, apol := range apols {
		found = true
		if !document {
			appendToUserAssociationInfo(ich, aws.ToString(user.UserName), aws.ToString(user.Arn), policyTypeAttached.String(), aws.ToString(apol.PolicyName), "")
			continue
		}
		doc, err := client.GetCustomerPolicyDocument(ctx, apol.PolicyArn, pols)
		if err != nil {
			return false, err
		}
		if len(filters) == 0 || api.Contains(doc, filters) {
			appendToUserAssociationInfo(ich, aws.ToString(user.UserName), aws.ToString(user.Arn), policyTypeAttached.String(), aws.ToString(apol.PolicyName), doc)
		}
	}
	return found, nil
}

func appendToUserAssociationInfoForAttachedGroupPolicy(ctx context.Context, client IIamClient, ich chan<- UserAssociationInfo, user types.User, group types.Group, document bool, filters []string, pols map[string]types.Policy) (bool, error) {
	found := false
	apols, err := client.GetAttachedGroupPolicies(ctx, group.GroupName)
	if err != nil {
		return false, err
	}
	for _, apol := range apols {
		found = true
		if !document {
			appendToUserAssociationInfo(ich, aws.ToString(user.UserName), aws.ToString(group.Arn), policyTypeAttached.String(), aws.ToString(apol.PolicyName), "")
			continue
		}
		doc, err := client.GetCustomerPolicyDocument(ctx, apol.PolicyArn, pols)
		if err != nil {
			return false, err
		}
		if len(filters) == 0 || api.Contains(doc, filters) {
			appendToUserAssociationInfo(ich, aws.ToString(user.UserName), aws.ToString(group.Arn), policyTypeAttached.String(), aws.ToString(apol.PolicyName), doc)
		}
	}
	return found, nil
}

func appendToUserAssociationInfoForInlineUserPolicy(ctx context.Context, client IIamClient, ich chan<- UserAssociationInfo, user types.User, document bool, filters []string) (bool, error) {
	found := false
	ipols, err := client.GetInlineUserPolicies(ctx, user.UserName)
	if err != nil {
		return false, err
	}
	for _, ipol := range ipols {
		found = true
		if !document {
			appendToUserAssociationInfo(ich, aws.ToString(user.UserName), aws.ToString(user.Arn), policyTypeInline.String(), ipol, "")
			continue
		}
		doc, err := client.GetInlineUserPolicyDocument(ctx, user.UserName, aws.String(ipol))
		if err != nil {
			return false, err
		}
		if len(filters) == 0 || api.Contains(doc, filters) {
			appendToUserAssociationInfo(ich, aws.ToString(user.UserName), aws.ToString(user.Arn), policyTypeInline.String(), ipol, doc)
		}
	}
	return found, nil
}

func appendToUserAssociationInfoForInlineGroupPolicy(ctx context.Context, client IIamClient, ich chan<- UserAssociationInfo, user types.User, group types.Group, document bool, filters []string) (bool, error) {
	found := false
	ipols, err := client.GetInlineGroupPolicies(ctx, group.GroupName)
	if err != nil {
		return false, err
	}
	for _, ipol := range ipols {
		found = true
		if !document {
			appendToUserAssociationInfo(ich, aws.ToString(user.UserName), aws.ToString(group.Arn), policyTypeInline.String(), ipol, "")
			continue
		}
		doc, err := client.GetInlineGroupPolicyDocument(ctx, group.GroupName, aws.String(ipol))
		if err != nil {
			return false, err
		}
		if len(filters) == 0 || api.Contains(doc, filters) {
			appendToUserAssociationInfo(ich, aws.ToString(user.UserName), aws.ToString(group.Arn), policyTypeInline.String(), ipol, doc)
		}
	}
	return found, nil
}

func appendToUserAssociationInfo(ich chan<- UserAssociationInfo, userName, attachedBy, policyType, policyName, policyDocument string) {
	ich <- UserAssociationInfo{
		UserName:       userName,
		AttachedBy:     attachedBy,
		PolicyType:     policyType,
		PolicyName:     policyName,
		PolicyDocument: policyDocument,
	}
}

package elbv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

var _ IElbClient = (*ElbClient)(nil)

type IElbClient interface {
	FetchTargetGroups(ctx context.Context, region string, targetType types.TargetTypeEnum) (map[string]types.TargetGroup, error)
	FetchTargets(ctx context.Context, region string, targetType types.TargetTypeEnum, reverse bool) ([]string, map[string][]string, error)
}

type ElbClient struct {
	*elasticloadbalancingv2.Client
}

func NewElbClient(cfg *aws.Config) *ElbClient {
	return &ElbClient{Client: elasticloadbalancingv2.NewFromConfig(*cfg)}
}

func (client *ElbClient) FetchTargetGroups(ctx context.Context, region string, targetType types.TargetTypeEnum) (map[string]types.TargetGroup, error) {
	return fetchElbv2TargetGroups(ctx, client.Client, region, targetType)
}

func (client *ElbClient) FetchTargets(ctx context.Context, region string, targetType types.TargetTypeEnum, reverse bool) ([]string, map[string][]string, error) {
	return fetchElbv2Targets(ctx, client.Client, region, targetType, reverse)
}

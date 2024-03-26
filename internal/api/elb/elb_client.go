package elb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
)

var _ IElbClient = (*ElbClient)(nil)

type IElbClient interface {
	FetchLoadBalancers(ctx context.Context, region string) (map[string]types.LoadBalancerDescription, error)
	FetchTargets(ctx context.Context, region string, reverse bool) ([]string, map[string][]string, error)
}

type ElbClient struct {
	*elasticloadbalancing.Client
}

func NewElbClient(cfg *aws.Config) *ElbClient {
	return &ElbClient{Client: elasticloadbalancing.NewFromConfig(*cfg)}
}

func (client *ElbClient) FetchLoadBalancers(ctx context.Context, region string) (map[string]types.LoadBalancerDescription, error) {
	return fetchElbLoadBalancers(ctx, client.Client, region)
}

func (client *ElbClient) FetchTargets(ctx context.Context, region string, reverse bool) ([]string, map[string][]string, error) {
	return fetchElbTargets(ctx, client.Client, region, reverse)
}

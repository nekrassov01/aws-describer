package elb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
)

func fetchElbLoadBalancers(ctx context.Context, client *elasticloadbalancing.Client, region string) (map[string]types.LoadBalancerDescription, error) {
	var marker *string
	res := make(map[string]types.LoadBalancerDescription)
	for {
		input := &elasticloadbalancing.DescribeLoadBalancersInput{
			Marker: marker,
		}
		opt := func(opt *elasticloadbalancing.Options) {
			opt.Region = region
		}
		o, err := client.DescribeLoadBalancers(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, lb := range o.LoadBalancerDescriptions {
			res[aws.ToString(lb.LoadBalancerName)] = lb
		}
		marker = o.NextMarker
		if marker == nil {
			break
		}
	}
	return res, nil
}

func fetchElbTargets(ctx context.Context, client *elasticloadbalancing.Client, region string, reverse bool) ([]string, map[string][]string, error) {
	var marker *string
	var lbs []types.LoadBalancerDescription
	for {
		input := &elasticloadbalancing.DescribeLoadBalancersInput{
			Marker: marker,
		}
		opt := func(opt *elasticloadbalancing.Options) {
			opt.Region = region
		}
		o, err := client.DescribeLoadBalancers(ctx, input, opt)
		if err != nil {
			return nil, nil, err
		}
		lbs = append(lbs, o.LoadBalancerDescriptions...)
		marker = o.NextMarker
		if marker == nil {
			break
		}
	}
	var ids []string
	idm := make(map[string][]string)
	for _, lb := range lbs {
		lbName := aws.ToString(lb.LoadBalancerName)
		for _, i := range lb.Instances {
			id := aws.ToString(i.InstanceId)
			key, value := lbName, id
			if reverse {
				key, value = id, lbName
			}
			idm[key] = append(idm[key], value)
		}
	}
	for id := range idm {
		ids = append(ids, id)
	}
	return ids, idm, nil
}

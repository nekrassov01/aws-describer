package elbv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

func fetchElbv2TargetGroups(ctx context.Context, client *elasticloadbalancingv2.Client, region string, targetType types.TargetTypeEnum) (map[string]types.TargetGroup, error) {
	var marker *string
	res := make(map[string]types.TargetGroup)
	for {
		input := &elasticloadbalancingv2.DescribeTargetGroupsInput{
			Marker: marker,
		}
		opt := func(opt *elasticloadbalancingv2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeTargetGroups(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, tg := range o.TargetGroups {
			if targetType == "" || tg.TargetType == targetType {
				res[aws.ToString(tg.TargetGroupName)] = tg
			}
		}
		marker = o.NextMarker
		if marker == nil {
			break
		}
	}
	return res, nil
}

func fetchElbv2Targets(ctx context.Context, client *elasticloadbalancingv2.Client, region string, targetType types.TargetTypeEnum, reverse bool) ([]string, map[string][]string, error) {
	var ids []string
	idm := make(map[string][]string)
	tgs, err := fetchElbv2TargetGroups(ctx, client, region, targetType)
	if err != nil {
		return nil, nil, err
	}
	for _, tg := range tgs {
		input := &elasticloadbalancingv2.DescribeTargetHealthInput{
			TargetGroupArn: tg.TargetGroupArn,
		}
		opt := func(opt *elasticloadbalancingv2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeTargetHealth(ctx, input, opt)
		if err != nil {
			return nil, nil, err
		}
		for _, thd := range o.TargetHealthDescriptions {
			id := aws.ToString(thd.Target.Id)
			tgName := aws.ToString(tg.TargetGroupName)
			key, value := tgName, id
			if reverse {
				key, value = id, tgName
			}
			idm[key] = append(idm[key], value)
		}
	}
	for id := range idm {
		ids = append(ids, id)
	}
	return ids, idm, nil
}

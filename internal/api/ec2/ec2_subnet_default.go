// Code generated by api/ec2/ec2_gen.go. DO NOT EDIT.

package ec2

import (
	"context"
	"runtime"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func DescribeSubnetInfo(ctx context.Context, cfg *aws.Config, regions []string, ids, names []string, filters []types.Filter, defaultFilter bool) ([]SubnetInfo, error) {
	client := NewEc2Client(cfg)
	eg, ctx := errgroup.WithContext(ctx)
	ich := make(chan SubnetInfo, runtime.NumCPU())
	var info []SubnetInfo
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range ich {
			info = append(info, i)
		}
	}()
	for _, region := range regions {
		region := region
		// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/throttling.html
		l := rate.NewLimiter(rate.Limit(50), 1)
		eg.Go(func() error {
			vpcs, err := client.FetchVpcs(ctx, region)
			if err != nil {
				return err
			}
			var token *string
			for {
				if err := l.Wait(ctx); err != nil {
					return err
				}
				input := CreateDescribeSubnetsInput(ids, names, filters, defaultFilter)
				input.NextToken = token
				opt := func(opt *ec2.Options) {
					opt.Region = region
				}
				o, err := client.DescribeSubnets(ctx, input, opt)
				if err != nil {
					return err
				}
				if err := GetSubnetInfo(ich, o.Subnets, region, vpcs); err != nil {
					return err
				}
				token = o.NextToken
				if token == nil {
					break
				}
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		close(ich)
		return nil, err
	}
	close(ich)
	wg.Wait()
	return info, nil
}

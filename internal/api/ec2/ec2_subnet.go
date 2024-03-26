package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func CreateDescribeSubnetsInput(ids, names []string, filters []types.Filter, defaultFilter bool) *ec2.DescribeSubnetsInput {
	var f []types.Filter
	if len(ids) > 0 {
		f = append(f, types.Filter{
			Name:   aws.String("subnet-id"),
			Values: ids,
		})
	}
	if len(names) > 0 {
		f = append(f, types.Filter{
			Name:   aws.String("tag:Name"),
			Values: names,
		})
	}
	if len(filters) > 0 {
		f = append(f, filters...)
	}
	if defaultFilter {
		f = append(f, types.Filter{
			Name:   aws.String("default-for-az"),
			Values: []string{"false"},
		})
	}
	return &ec2.DescribeSubnetsInput{
		Filters: f,
	}
}

type SubnetInfo struct {
	SubnetId                string
	SubnetName              string
	AvailabilityZone        string
	AvailableIpAddressCount int32
	DefaultForAz            bool
	State                   types.SubnetState
	VpcId                   string
	VpcName                 string
	AddressType             string
	CidrBlock               string
	Region                  string
}

func GetSubnetInfo(ich chan<- SubnetInfo, subnets []types.Subnet, region string, vpcs map[string]types.Vpc) error {
	for _, subnet := range subnets {
		vpcId := aws.ToString(subnet.VpcId)
		vpc, err := findEc2VpcById(vpcId, vpcs)
		if err != nil {
			return err
		}
		obj := SubnetInfo{
			SubnetId:                aws.ToString(subnet.SubnetId),
			SubnetName:              getEc2NameTagValue(subnet.Tags),
			AvailabilityZone:        aws.ToString(subnet.AvailabilityZone),
			AvailableIpAddressCount: aws.ToInt32(subnet.AvailableIpAddressCount),
			DefaultForAz:            aws.ToBool(subnet.DefaultForAz),
			State:                   subnet.State,
			VpcId:                   vpcId,
			VpcName:                 getEc2NameTagValue(vpc.Tags),
			AddressType:             addressTypeIpv4.String(),
			CidrBlock:               aws.ToString(subnet.CidrBlock),
			Region:                  region,
		}
		ich <- obj
		for _, assoc := range subnet.Ipv6CidrBlockAssociationSet {
			obj.AddressType = addressTypeIpv6.String()
			obj.CidrBlock = aws.ToString(assoc.Ipv6CidrBlock)
			ich <- obj
		}
	}
	return nil
}

type SubnetRouteInfo struct {
	SubnetId         string
	SubnetName       string
	AvailabilityZone string
	VpcId            string
	VpcName          string
	RouteTableId     string
	RouteTableName   string
	DestinationType  string
	Destination      string
	TargetType       string
	Target           string
	Region           string
}

func FetchDataForSubnetRouteInfo(ctx context.Context, l *rate.Limiter, client IEc2Client, region string) (map[string]types.Vpc, map[string]types.RouteTable, error) {
	vpcs := make(map[string]types.Vpc)
	rtbs := make(map[string]types.RouteTable)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		vpcs, err = client.FetchVpcs(ctx, region)
		return err
	})
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		rtbs, err = client.FetchRouteTables(ctx, region)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}
	return vpcs, rtbs, nil
}

func GetSubnetRouteInfo(ich chan<- SubnetRouteInfo, subnets []types.Subnet, region string, vpcs map[string]types.Vpc, rtbs map[string]types.RouteTable) error {
	for _, subnet := range subnets {
		vpcId := aws.ToString(subnet.VpcId)
		vpc, err := findEc2VpcById(vpcId, vpcs)
		if err != nil {
			return err
		}
		rtb, err := findEc2RouteTableBySubnet(subnet, rtbs)
		if err != nil {
			return err
		}
		rtbId := aws.ToString(rtb.RouteTableId)
		obj := SubnetRouteInfo{
			SubnetId:         aws.ToString(subnet.SubnetId),
			SubnetName:       getEc2NameTagValue(subnet.Tags),
			AvailabilityZone: aws.ToString(subnet.AvailabilityZone),
			VpcId:            vpcId,
			VpcName:          getEc2NameTagValue(vpc.Tags),
			RouteTableId:     rtbId,
			RouteTableName:   getEc2NameTagValue(rtb.Tags),
			Region:           region,
		}
		for _, rt := range rtb.Routes {
			destinationType, destination, err := getEc2RouteDestination(rt)
			if err != nil {
				return err
			}
			targetType, target, err := getEc2RouteTarget(rt)
			if err != nil {
				return err
			}
			obj.DestinationType = destinationType
			obj.Destination = destination
			obj.TargetType = targetType
			obj.Target = target
			ich <- obj
		}
	}
	return nil
}

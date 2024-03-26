package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func CreateDescribeRouteTablesInput(ids, names []string, filters []types.Filter, defaultFilter bool) *ec2.DescribeRouteTablesInput {
	var f []types.Filter
	if len(ids) > 0 {
		f = append(f, types.Filter{
			Name:   aws.String("route-table-id"),
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
			Name:   aws.String("association.main"),
			Values: []string{"false"},
		})
	}
	return &ec2.DescribeRouteTablesInput{
		Filters: f,
	}
}

type RouteTableInfo struct {
	RouteTableId    string
	RouteTableName  string
	VpcId           string
	VpcName         string
	DestinationType string
	Destination     string
	TargetType      string
	Target          string
	State           types.RouteState
	Region          string
}

func GetRouteTableInfo(ich chan<- RouteTableInfo, rtbs []types.RouteTable, region string, vpcs map[string]types.Vpc) error {
	for _, rtb := range rtbs {
		routes, err := handleRoutes(rtb, vpcs, region)
		if err != nil {
			return err
		}
		for _, route := range routes {
			ich <- route
		}
	}
	return nil
}

func handleRoutes(rtb types.RouteTable, vpcs map[string]types.Vpc, region string) ([]RouteTableInfo, error) {
	var info []RouteTableInfo
	vpcId := aws.ToString(rtb.VpcId)
	vpc, err := findEc2VpcById(vpcId, vpcs)
	if err != nil {
		return nil, err
	}
	for _, rt := range rtb.Routes {
		destinationType, destination, err := getEc2RouteDestination(rt)
		if err != nil {
			return nil, err
		}
		targetType, target, err := getEc2RouteTarget(rt)
		if err != nil {
			return nil, err
		}
		info = append(info, RouteTableInfo{
			RouteTableId:    aws.ToString(rtb.RouteTableId),
			RouteTableName:  getEc2NameTagValue(rtb.Tags),
			VpcId:           vpcId,
			VpcName:         getEc2NameTagValue(vpc.Tags),
			DestinationType: destinationType,
			Destination:     destination,
			TargetType:      targetType,
			Target:          target,
			State:           rt.State,
			Region:          region,
		})
	}
	return info, nil
}

type RouteTableAssociationInfo struct {
	RouteTableId   string
	RouteTableName string
	VpcId          string
	VpcName        string
	Main           bool
	SubnetId       string
	SubnetName     string
	State          types.RouteTableAssociationStateCode
	Region         string
}

func FetchDataForRouteTableAssociationInfo(ctx context.Context, l *rate.Limiter, client IEc2Client, region string) (map[string]types.Vpc, map[string]types.Subnet, error) {
	vpcs := make(map[string]types.Vpc)
	sbns := make(map[string]types.Subnet)
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
		sbns, err = client.FetchSubnets(ctx, region)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}
	return vpcs, sbns, nil
}

func GetRouteTableAssociationInfo(ich chan<- RouteTableAssociationInfo, rtbs []types.RouteTable, region string, vpcs map[string]types.Vpc, sbns map[string]types.Subnet) error {
	for _, rtb := range rtbs {
		vpcId := aws.ToString(rtb.VpcId)
		vpc, err := findEc2VpcById(vpcId, vpcs)
		if err != nil {
			return err
		}
		for _, assoc := range rtb.Associations {
			subnetId := aws.ToString(assoc.SubnetId)
			subnetName := ""
			if subnetId != "" {
				sbn, err := findEc2SubnetById(subnetId, sbns)
				if err != nil {
					return err
				}
				subnetName = getEc2NameTagValue(sbn.Tags)
			}
			ich <- RouteTableAssociationInfo{
				RouteTableId:   aws.ToString(rtb.RouteTableId),
				RouteTableName: getEc2NameTagValue(rtb.Tags),
				VpcId:          vpcId,
				VpcName:        getEc2NameTagValue(vpc.Tags),
				Main:           aws.ToBool(assoc.Main),
				SubnetId:       subnetId,
				SubnetName:     subnetName,
				State:          assoc.AssociationState.State,
				Region:         region,
			}
		}
	}
	return nil
}

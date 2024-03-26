package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func CreateDescribeVpcsInput(ids, names []string, filters []types.Filter, defaultFilter bool) *ec2.DescribeVpcsInput {
	var f []types.Filter
	if len(ids) > 0 {
		f = append(f, types.Filter{
			Name:   aws.String("vpc-id"),
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
			Name:   aws.String("is-default"),
			Values: []string{"false"},
		})
	}
	return &ec2.DescribeVpcsInput{
		Filters: f,
	}
}

type VpcInfo struct {
	VpcId           string
	VpcName         string
	DhcpOptionsId   string
	DhcpOptionsName string
	IsDefault       bool
	InstanceTenancy types.Tenancy
	OwnerId         string
	Region          string
}

func GetVpcInfo(ich chan<- VpcInfo, vpcs []types.Vpc, region string, dopts map[string]types.DhcpOptions) error {
	for _, vpc := range vpcs {
		dhcpOptId := aws.ToString(vpc.DhcpOptionsId)
		dhcpOpt, err := findEc2DhcpOptionById(dhcpOptId, dopts)
		if err != nil {
			return err
		}
		ich <- VpcInfo{
			VpcId:           aws.ToString(vpc.VpcId),
			VpcName:         getEc2NameTagValue(vpc.Tags),
			DhcpOptionsId:   dhcpOptId,
			DhcpOptionsName: getEc2NameTagValue(dhcpOpt.Tags),
			IsDefault:       aws.ToBool(vpc.IsDefault),
			InstanceTenancy: vpc.InstanceTenancy,
			OwnerId:         aws.ToString(vpc.OwnerId),
			Region:          region,
		}
	}
	return nil
}

type VpcAttributeInfo struct {
	VpcId              string
	VpcName            string
	EnableDnsSupport   bool
	EnableDnsHostnames bool
	DhcpOptionsId      string
	DhcpOptionsName    string
	IsDefault          bool
	InstanceTenancy    types.Tenancy
	OwnerId            string
	Region             string
}

func GetVpcAttributeInfo(ctx context.Context, l *rate.Limiter, client IEc2Client, ich chan<- VpcAttributeInfo, vpcs []types.Vpc, region string, dopts map[string]types.DhcpOptions) error {
	for _, vpc := range vpcs {
		var enableDnsSupport, enableDnsHostnames bool
		eg, ctx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			if err := l.Wait(ctx); err != nil {
				return err
			}
			var err error
			enableDnsSupport, err = client.GetVpcDnsSupport(ctx, region, vpc.VpcId)
			return err
		})
		eg.Go(func() error {
			if err := l.Wait(ctx); err != nil {
				return err
			}
			var err error
			enableDnsHostnames, err = client.GetVpcDnsHostnames(ctx, region, vpc.VpcId)
			return err
		})
		if err := eg.Wait(); err != nil {
			return err
		}
		doptId := aws.ToString(vpc.DhcpOptionsId)
		dopt, err := findEc2DhcpOptionById(doptId, dopts)
		if err != nil {
			return err
		}
		ich <- VpcAttributeInfo{
			VpcId:              aws.ToString(vpc.VpcId),
			VpcName:            getEc2NameTagValue(vpc.Tags),
			EnableDnsSupport:   enableDnsSupport,
			EnableDnsHostnames: enableDnsHostnames,
			DhcpOptionsId:      doptId,
			DhcpOptionsName:    getEc2NameTagValue(dopt.Tags),
			IsDefault:          aws.ToBool(vpc.IsDefault),
			InstanceTenancy:    vpc.InstanceTenancy,
			OwnerId:            aws.ToString(vpc.OwnerId),
			Region:             region,
		}
	}
	return nil
}

type VpcCidrInfo struct {
	VpcId              string
	VpcName            string
	DhcpOptionsId      string
	DhcpOptionsName    string
	IsDefault          bool
	InstanceTenancy    types.Tenancy
	OwnerId            string
	State              types.VpcState
	AddressType        string
	CidrBlock          string
	NetworkBorderGroup string
	Pool               string
	Region             string
}

func GetVpcCidrInfo(ich chan<- VpcCidrInfo, vpcs []types.Vpc, region string, dopts map[string]types.DhcpOptions) error {
	for _, vpc := range vpcs {
		dhcpOptId := aws.ToString(vpc.DhcpOptionsId)
		dhcpOpt, err := findEc2DhcpOptionById(dhcpOptId, dopts)
		if err != nil {
			return err
		}
		obj := VpcCidrInfo{
			VpcId:           aws.ToString(vpc.VpcId),
			VpcName:         getEc2NameTagValue(vpc.Tags),
			DhcpOptionsId:   dhcpOptId,
			DhcpOptionsName: getEc2NameTagValue(dhcpOpt.Tags),
			IsDefault:       aws.ToBool(vpc.IsDefault),
			InstanceTenancy: vpc.InstanceTenancy,
			OwnerId:         aws.ToString(vpc.OwnerId),
			State:           vpc.State,
			Region:          region,
		}
		for _, assoc := range vpc.CidrBlockAssociationSet {
			obj.AddressType = addressTypeIpv4.String()
			obj.CidrBlock = aws.ToString(assoc.CidrBlock)
			ich <- obj
		}
		for _, assoc := range vpc.Ipv6CidrBlockAssociationSet {
			obj.AddressType = addressTypeIpv6.String()
			obj.CidrBlock = aws.ToString(assoc.Ipv6CidrBlock)
			obj.NetworkBorderGroup = aws.ToString(assoc.NetworkBorderGroup)
			obj.Pool = aws.ToString(assoc.Ipv6Pool)
			ich <- obj
		}
	}
	return nil
}

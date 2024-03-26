package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var _ IEc2Client = (*Ec2Client)(nil)

type IEc2Client interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	DescribeImages(ctx context.Context, params *ec2.DescribeImagesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error)
	DescribeSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
	DescribeVpcs(ctx context.Context, params *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error)
	DescribeSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
	DescribeRouteTables(ctx context.Context, params *ec2.DescribeRouteTablesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error)

	FetchImages(ctx context.Context, region string) (map[string]types.Image, error)
	FetchSnapshots(ctx context.Context, region string) (map[string]types.Snapshot, error)
	FetchVolumes(ctx context.Context, region string) (map[string]types.Volume, error)
	FetchSecurityGroups(ctx context.Context, region string) (map[string]types.SecurityGroup, error)
	FetchVpcs(ctx context.Context, region string) (map[string]types.Vpc, error)
	FetchSubnets(ctx context.Context, region string) (map[string]types.Subnet, error)
	FetchRouteTables(ctx context.Context, region string) (map[string]types.RouteTable, error)
	FetchPrefixLists(ctx context.Context, region string) (map[string]types.PrefixList, error)
	FetchManagedPrefixLists(ctx context.Context, region string) (map[string]types.ManagedPrefixList, error)
	FetchDhcpOptions(ctx context.Context, region string) (map[string]types.DhcpOptions, error)
	GetVpcDnsSupport(ctx context.Context, region string, id *string) (bool, error)
	GetVpcDnsHostnames(ctx context.Context, region string, id *string) (bool, error)
}

type Ec2Client struct {
	*ec2.Client
}

func NewEc2Client(cfg *aws.Config) *Ec2Client {
	return &Ec2Client{Client: ec2.NewFromConfig(*cfg)}
}

func (client *Ec2Client) FetchImages(ctx context.Context, region string) (map[string]types.Image, error) {
	return fetchEc2Images(ctx, client.Client, region)
}

func (client *Ec2Client) FetchSnapshots(ctx context.Context, region string) (map[string]types.Snapshot, error) {
	return fetchEc2Snapshots(ctx, client.Client, region)
}

func (client *Ec2Client) FetchVolumes(ctx context.Context, region string) (map[string]types.Volume, error) {
	return fetchEc2Volumes(ctx, client.Client, region)
}

func (client *Ec2Client) FetchSecurityGroups(ctx context.Context, region string) (map[string]types.SecurityGroup, error) {
	return fetchEc2SecurityGroups(ctx, client.Client, region)
}

func (client *Ec2Client) FetchVpcs(ctx context.Context, region string) (map[string]types.Vpc, error) {
	return fetchEc2Vpcs(ctx, client.Client, region)
}

func (client *Ec2Client) FetchSubnets(ctx context.Context, region string) (map[string]types.Subnet, error) {
	return fetchEc2Subnets(ctx, client.Client, region)
}

func (client *Ec2Client) FetchRouteTables(ctx context.Context, region string) (map[string]types.RouteTable, error) {
	return fetchEc2RouteTables(ctx, client.Client, region)
}

func (client *Ec2Client) FetchPrefixLists(ctx context.Context, region string) (map[string]types.PrefixList, error) {
	return fetchEc2PrefixLists(ctx, client.Client, region)
}

func (client *Ec2Client) FetchManagedPrefixLists(ctx context.Context, region string) (map[string]types.ManagedPrefixList, error) {
	return fetchEc2ManagedPrefixLists(ctx, client.Client, region)
}

func (client *Ec2Client) FetchDhcpOptions(ctx context.Context, region string) (map[string]types.DhcpOptions, error) {
	return fetchEc2DhcpOptions(ctx, client.Client, region)
}

func (client *Ec2Client) GetVpcDnsSupport(ctx context.Context, region string, id *string) (bool, error) {
	return getEc2VpcDnsSupport(ctx, client.Client, region, id)
}

func (client *Ec2Client) GetVpcDnsHostnames(ctx context.Context, region string, id *string) (bool, error) {
	return getEc2VpcDnsHostnames(ctx, client.Client, region, id)
}

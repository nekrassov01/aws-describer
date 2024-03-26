package ec2

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var instanceValidState = []string{
	"pending",
	"running",
	"stopping",
	"stopped",
}

func fetchEc2Instances(ctx context.Context, client *ec2.Client, region string) (map[string]types.Instance, error) {
	var token *string
	res := make(map[string]types.Instance)
	for {
		input := &ec2.DescribeInstancesInput{
			NextToken: token,
			Filters: []types.Filter{
				{
					Name:   aws.String("instance-state-name"),
					Values: instanceValidState,
				},
			},
		}
		opt := func(opt *ec2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeInstances(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, r := range o.Reservations {
			for _, i := range r.Instances {
				res[aws.ToString(i.InstanceId)] = i
			}
		}
		token = o.NextToken
		if token == nil {
			break
		}
	}
	return res, nil
}

func fetchEc2Images(ctx context.Context, client *ec2.Client, region string) (map[string]types.Image, error) {
	instances, err := fetchEc2Instances(ctx, client, region)
	if err != nil {
		return nil, err
	}
	var imgIds []string
	for _, instance := range instances {
		imgIds = append(imgIds, aws.ToString(instance.ImageId))
	}
	if imgIds == nil {
		return nil, nil
	}
	slices.Sort(imgIds)

	var token *string
	res := make(map[string]types.Image)
	for {
		input := &ec2.DescribeImagesInput{
			NextToken: token,
			ImageIds:  slices.Compact(imgIds),
		}
		opt := func(opt *ec2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeImages(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, i := range o.Images {
			res[aws.ToString(i.ImageId)] = i
		}
		token = o.NextToken
		if token == nil {
			break
		}
	}
	return res, nil
}

func fetchEc2Snapshots(ctx context.Context, client *ec2.Client, region string) (map[string]types.Snapshot, error) {
	var token *string
	res := make(map[string]types.Snapshot)
	for {
		input := &ec2.DescribeSnapshotsInput{
			NextToken: token,
			OwnerIds:  []string{"self"},
		}
		opt := func(opt *ec2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeSnapshots(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, s := range o.Snapshots {
			res[aws.ToString(s.SnapshotId)] = s
		}
		token = o.NextToken
		if token == nil {
			break
		}
	}
	return res, nil
}

func fetchEc2Volumes(ctx context.Context, client *ec2.Client, region string) (map[string]types.Volume, error) {
	var token *string
	res := make(map[string]types.Volume)
	for {
		input := &ec2.DescribeVolumesInput{
			NextToken: token,
		}
		opt := func(opt *ec2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeVolumes(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, v := range o.Volumes {
			res[aws.ToString(v.VolumeId)] = v
		}
		token = o.NextToken
		if token == nil {
			break
		}
	}
	return res, nil
}

func fetchEc2SecurityGroups(ctx context.Context, client *ec2.Client, region string) (map[string]types.SecurityGroup, error) {
	var token *string
	res := make(map[string]types.SecurityGroup)
	for {
		input := &ec2.DescribeSecurityGroupsInput{
			NextToken: token,
		}
		opt := func(opt *ec2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeSecurityGroups(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, sg := range o.SecurityGroups {
			res[aws.ToString(sg.GroupId)] = sg
		}
		token = o.NextToken
		if token == nil {
			break
		}
	}
	return res, nil
}

func fetchEc2Vpcs(ctx context.Context, client *ec2.Client, region string) (map[string]types.Vpc, error) {
	var token *string
	res := make(map[string]types.Vpc)
	for {
		input := &ec2.DescribeVpcsInput{
			NextToken: token,
		}
		opt := func(opt *ec2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeVpcs(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, vpc := range o.Vpcs {
			res[aws.ToString(vpc.VpcId)] = vpc
		}
		token = o.NextToken
		if token == nil {
			break
		}
	}
	return res, nil
}

func fetchEc2Subnets(ctx context.Context, client *ec2.Client, region string) (map[string]types.Subnet, error) {
	var token *string
	res := make(map[string]types.Subnet)
	for {
		input := &ec2.DescribeSubnetsInput{
			NextToken: token,
		}
		opt := func(opt *ec2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeSubnets(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, subnet := range o.Subnets {
			res[aws.ToString(subnet.SubnetId)] = subnet
		}
		token = o.NextToken
		if token == nil {
			break
		}
	}
	return res, nil
}

func fetchEc2RouteTables(ctx context.Context, client *ec2.Client, region string) (map[string]types.RouteTable, error) {
	var token *string
	res := make(map[string]types.RouteTable)
	for {
		input := &ec2.DescribeRouteTablesInput{
			NextToken: token,
		}
		opt := func(opt *ec2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeRouteTables(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, rtb := range o.RouteTables {
			res[aws.ToString(rtb.RouteTableId)] = rtb
		}
		token = o.NextToken
		if token == nil {
			break
		}
	}
	return res, nil
}

func fetchEc2PrefixLists(ctx context.Context, client *ec2.Client, region string) (map[string]types.PrefixList, error) {
	var token *string
	res := make(map[string]types.PrefixList)
	for {
		input := &ec2.DescribePrefixListsInput{
			NextToken: token,
		}
		opt := func(opt *ec2.Options) {
			opt.Region = region
		}
		o, err := client.DescribePrefixLists(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, vpc := range o.PrefixLists {
			res[aws.ToString(vpc.PrefixListId)] = vpc
		}
		token = o.NextToken
		if token == nil {
			break
		}
	}
	return res, nil
}

func fetchEc2ManagedPrefixLists(ctx context.Context, client *ec2.Client, region string) (map[string]types.ManagedPrefixList, error) {
	var token *string
	res := make(map[string]types.ManagedPrefixList)
	for {
		input := &ec2.DescribeManagedPrefixListsInput{
			NextToken: token,
		}
		opt := func(opt *ec2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeManagedPrefixLists(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, vpc := range o.PrefixLists {
			res[aws.ToString(vpc.PrefixListId)] = vpc
		}
		token = o.NextToken
		if token == nil {
			break
		}
	}
	return res, nil
}

func fetchEc2DhcpOptions(ctx context.Context, client *ec2.Client, region string) (map[string]types.DhcpOptions, error) {
	var token *string
	res := make(map[string]types.DhcpOptions)
	for {
		input := &ec2.DescribeDhcpOptionsInput{
			NextToken: token,
		}
		opt := func(opt *ec2.Options) {
			opt.Region = region
		}
		o, err := client.DescribeDhcpOptions(ctx, input, opt)
		if err != nil {
			return nil, err
		}
		for _, dopt := range o.DhcpOptions {
			res[aws.ToString(dopt.DhcpOptionsId)] = dopt
		}
		token = o.NextToken
		if token == nil {
			break
		}
	}
	return res, nil
}

func getEc2VpcDnsSupport(ctx context.Context, client *ec2.Client, region string, id *string) (bool, error) {
	input := &ec2.DescribeVpcAttributeInput{
		VpcId:     id,
		Attribute: types.VpcAttributeNameEnableDnsSupport,
	}
	opt := func(opt *ec2.Options) {
		opt.Region = region
	}
	attr, err := client.DescribeVpcAttribute(ctx, input, opt)
	if err != nil {
		return false, err
	}
	return aws.ToBool(attr.EnableDnsSupport.Value), err
}

func getEc2VpcDnsHostnames(ctx context.Context, client *ec2.Client, region string, id *string) (bool, error) {
	input := &ec2.DescribeVpcAttributeInput{
		VpcId:     id,
		Attribute: types.VpcAttributeNameEnableDnsHostnames,
	}
	opt := func(opt *ec2.Options) {
		opt.Region = region
	}
	attr, err := client.DescribeVpcAttribute(ctx, input, opt)
	if err != nil {
		return false, err
	}
	return aws.ToBool(attr.EnableDnsHostnames.Value), err
}

func getEc2NameTagValue(tags []types.Tag) string {
	for _, t := range tags {
		if t.Key != nil && strings.EqualFold(aws.ToString(t.Key), "Name") && t.Value != nil {
			return *t.Value
		}
	}
	return ""
}

func findEc2ImageById(id string, m map[string]types.Image) *types.Image {
	if item, ok := m[id]; ok {
		return &item
	}
	return nil
}

func findEc2SnapshotById(id string, m map[string]types.Snapshot) *types.Snapshot {
	if item, ok := m[id]; ok {
		return &item
	}
	return nil
}

func findEc2VolumeById(id string, m map[string]types.Volume) *types.Volume {
	if item, ok := m[id]; ok {
		return &item
	}
	return nil
}

func findEc2SecurityGroupById(id string, m map[string]types.SecurityGroup) (*types.SecurityGroup, error) {
	if item, ok := m[id]; ok {
		return &item, nil
	}
	return nil, fmt.Errorf("no security group found: %s", id)
}

func findEc2VpcById(id string, m map[string]types.Vpc) (*types.Vpc, error) {
	if item, ok := m[id]; ok {
		return &item, nil
	}
	return nil, fmt.Errorf("no vpc found %s: ", id)
}

func findEc2SubnetById(id string, m map[string]types.Subnet) (*types.Subnet, error) {
	if item, ok := m[id]; ok {
		return &item, nil
	}
	return nil, fmt.Errorf("no subnet found %s: ", id)
}

func findEc2RouteTableBySubnet(subnet types.Subnet, m map[string]types.RouteTable) (*types.RouteTable, error) {
	if item, ok := m[aws.ToString(subnet.SubnetId)]; ok {
		return &item, nil
	}
	if subnet.VpcId != nil {
		for _, item := range m {
			if item.VpcId != nil && aws.ToString(item.VpcId) == aws.ToString(subnet.VpcId) {
				for _, assoc := range item.Associations {
					if aws.ToBool(assoc.Main) {
						return &item, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("no route table found %s: ", *subnet.SubnetId)
}

func findEc2PrefixListById(id string, m map[string]types.PrefixList) *types.PrefixList {
	if item, ok := m[id]; ok {
		return &item
	}
	return nil
}

func findEc2ManagedPrefixListById(id string, m map[string]types.ManagedPrefixList) *types.ManagedPrefixList {
	if item, ok := m[id]; ok {
		return &item
	}
	return nil
}

func findEc2PrefixListNameById(id string, upls map[string]types.PrefixList, mpls map[string]types.ManagedPrefixList) (string, error) {
	upl := findEc2PrefixListById(id, upls)
	if upl != nil && upl.PrefixListName != nil {
		return aws.ToString(upl.PrefixListName), nil
	}
	mpl := findEc2ManagedPrefixListById(id, mpls)
	if mpl != nil && mpl.PrefixListName != nil {
		return aws.ToString(mpl.PrefixListName), nil
	}
	return "", fmt.Errorf("no prefix list found %s: ", id)
}

func findEc2DhcpOptionById(id string, m map[string]types.DhcpOptions) (*types.DhcpOptions, error) {
	if item, ok := m[id]; ok {
		return &item, nil
	}
	return nil, fmt.Errorf("no dhcp options found %s: ", id)
}

func getEc2RouteDestination(rt types.Route) (string, string, error) {
	var dtype, d string
	switch {
	case rt.DestinationCidrBlock != nil:
		dtype = addressTypeIpv4.String()
		d = aws.ToString(rt.DestinationCidrBlock)
	case rt.DestinationIpv6CidrBlock != nil:
		dtype = addressTypeIpv6.String()
		d = aws.ToString(rt.DestinationIpv6CidrBlock)
	case rt.DestinationPrefixListId != nil:
		dtype = addressTypePrefixList.String()
		d = aws.ToString(rt.DestinationPrefixListId)
	default:
		return "", "", fmt.Errorf("unknown destination in route: valid values: %s", strings.Join(addressTypes, "|"))
	}
	return dtype, d, nil
}

func getEc2RouteTarget(rt types.Route) (string, string, error) {
	var ttype, t string
	switch {
	case rt.GatewayId != nil:
		t = aws.ToString(rt.GatewayId)
		switch {
		case strings.HasPrefix(t, "local"):
			ttype = targetTypeLocal.String()
		case strings.HasPrefix(t, "igw"):
			ttype = targetTypeInternetGateway.String()
		case strings.HasPrefix(t, "vgw"):
			ttype = targetTypeVpnGateway.String()
		case strings.HasPrefix(t, "vpce"):
			ttype = targetTypeVpcEndpoint.String()
		default:
			ttype = targetTypeOther.String()
		}
	case rt.NatGatewayId != nil:
		ttype = targetTypeNatGateway.String()
		t = aws.ToString(rt.NatGatewayId)
	case rt.VpcPeeringConnectionId != nil:
		ttype = targetTypeVpcPeeringConnection.String()
		t = aws.ToString(rt.VpcPeeringConnectionId)
	case rt.TransitGatewayId != nil:
		ttype = targetTypeTransitGateway.String()
		t = aws.ToString(rt.TransitGatewayId)
	case rt.EgressOnlyInternetGatewayId != nil:
		ttype = targetTypeEgressOnlyInternetGateway.String()
		t = aws.ToString(rt.EgressOnlyInternetGatewayId)
	case rt.CarrierGatewayId != nil:
		ttype = targetTypeCarrierGateway.String()
		t = aws.ToString(rt.CarrierGatewayId)
	case rt.InstanceId != nil:
		ttype = targetTypeInstance.String()
		t = aws.ToString(rt.InstanceId) + "/" + aws.ToString(rt.InstanceOwnerId)
	case rt.NetworkInterfaceId != nil:
		ttype = targetTypeNetworkInterface.String()
		t = aws.ToString(rt.NetworkInterfaceId)
	case rt.LocalGatewayId != nil:
		ttype = targetTypeOutpostLocalGateway.String()
		t = aws.ToString(rt.LocalGatewayId)
	case rt.CoreNetworkArn != nil:
		ttype = targetTypeCoreNetwork.String()
		t = aws.ToString(rt.CoreNetworkArn)
	default:
		return "", "", fmt.Errorf("unknown target in route: valid values: %s", strings.Join(targetTypes, "|"))
	}
	return ttype, t, nil
}

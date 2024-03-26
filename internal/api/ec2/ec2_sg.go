package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func CreateDescribeSecurityGroupsInput(ids, names []string, filters []types.Filter, _ bool) *ec2.DescribeSecurityGroupsInput {
	var f []types.Filter
	if ids != nil {
		f = append(f, types.Filter{
			Name:   aws.String("group-id"),
			Values: ids,
		})
	}
	if names != nil {
		f = append(f, types.Filter{
			Name:   aws.String("group-name"),
			Values: ids,
		})
	}
	if len(filters) > 0 {
		f = append(f, filters...)
	}
	return &ec2.DescribeSecurityGroupsInput{
		Filters: f,
	}
}

type SecurityGroupInfo struct {
	SecurityGroupId   string
	SecurityGroupName string
	VpcId             string
	VpcName           string
	Region            string
}

func GetSecurityGroupInfo(ich chan<- SecurityGroupInfo, sgs []types.SecurityGroup, region string, vpcs map[string]types.Vpc) error {
	for _, sg := range sgs {
		vpcId := aws.ToString(sg.VpcId)
		vpc, err := findEc2VpcById(vpcId, vpcs)
		if err != nil {
			return err
		}
		ich <- SecurityGroupInfo{
			SecurityGroupId:   aws.ToString(sg.GroupId),
			SecurityGroupName: aws.ToString(sg.GroupName),
			VpcId:             vpcId,
			VpcName:           getEc2NameTagValue(vpc.Tags),
			Region:            region,
		}
	}
	return nil
}

type SecurityGroupPermissionsInfo struct {
	SecurityGroupId   string
	SecurityGroupName string
	VpcId             string
	VpcName           string
	FlowDirection     string
	IpProtocol        string
	FromPort          int32
	ToPort            int32
	AddressType       string
	CidrBlock         string
	Region            string
}

func FetchDataForSecurityGroupPermissionsInfo(ctx context.Context, l *rate.Limiter, client IEc2Client, region string) (map[string]types.Vpc, map[string]types.PrefixList, map[string]types.ManagedPrefixList, error) {
	vpcs := make(map[string]types.Vpc)
	upls := make(map[string]types.PrefixList)
	mpls := make(map[string]types.ManagedPrefixList)
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
		upls, err = client.FetchPrefixLists(ctx, region)
		return err
	})
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		mpls, err = client.FetchManagedPrefixLists(ctx, region)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, nil, nil, err
	}
	return vpcs, upls, mpls, nil
}

func GetSecurityGroupPermissionsInfo(ich chan<- SecurityGroupPermissionsInfo, sgs []types.SecurityGroup, region string, vpcs map[string]types.Vpc, upls map[string]types.PrefixList, mpls map[string]types.ManagedPrefixList) error {
	for _, sg := range sgs {
		if err := appendToSecurityGroupPermissionsInfo(ich, sg, vpcs, upls, mpls, false, region); err != nil {
			return err
		}
		if err := appendToSecurityGroupPermissionsInfo(ich, sg, vpcs, upls, mpls, true, region); err != nil {
			return err
		}
	}
	return nil
}

func appendToSecurityGroupPermissionsInfo(ich chan<- SecurityGroupPermissionsInfo, item types.SecurityGroup, vpcs map[string]types.Vpc, upls map[string]types.PrefixList, mpls map[string]types.ManagedPrefixList, isEgress bool, region string) error {
	perms, err := handleSecurityGroupPermissionsInfo(item, vpcs, upls, mpls, isEgress, region)
	if err != nil {
		return err
	}
	for _, perm := range perms {
		ich <- perm
	}
	return nil
}

func handleSecurityGroupPermissionsInfo(item types.SecurityGroup, vpcs map[string]types.Vpc, upls map[string]types.PrefixList, mpls map[string]types.ManagedPrefixList, isEgress bool, region string) ([]SecurityGroupPermissionsInfo, error) {
	var ipPermissions []types.IpPermission
	var flowDirection string
	var info []SecurityGroupPermissionsInfo
	if isEgress {
		ipPermissions = item.IpPermissionsEgress
		flowDirection = Egress.String()
	} else {
		ipPermissions = item.IpPermissions
		flowDirection = Ingress.String()
	}
	vpcId := aws.ToString(item.VpcId)
	vpc, err := findEc2VpcById(vpcId, vpcs)
	if err != nil {
		return nil, err
	}
	for _, ipPermission := range ipPermissions {
		obj := SecurityGroupPermissionsInfo{
			SecurityGroupId:   aws.ToString(item.GroupId),
			SecurityGroupName: aws.ToString(item.GroupName),
			VpcId:             vpcId,
			VpcName:           getEc2NameTagValue(vpc.Tags),
			FlowDirection:     flowDirection,
			IpProtocol:        aws.ToString(ipPermission.IpProtocol),
			FromPort:          aws.ToInt32(ipPermission.FromPort),
			ToPort:            aws.ToInt32(ipPermission.ToPort),
			Region:            region,
		}
		for _, ipv4 := range ipPermission.IpRanges {
			obj.AddressType = addressTypeIpv4.String()
			obj.CidrBlock = aws.ToString(ipv4.CidrIp)
			info = append(info, obj)
		}
		for _, ipv6 := range ipPermission.Ipv6Ranges {
			obj.AddressType = addressTypeIpv6.String()
			obj.CidrBlock = aws.ToString(ipv6.CidrIpv6)
			info = append(info, obj)
		}
		for _, group := range ipPermission.UserIdGroupPairs {
			obj.AddressType = addressTypeSecurityGroup.String()
			obj.CidrBlock = aws.ToString(group.GroupId) + "/" + aws.ToString(group.UserId)
			info = append(info, obj)
		}
		for _, prefixList := range ipPermission.PrefixListIds {
			name, err := findEc2PrefixListNameById(aws.ToString(prefixList.PrefixListId), upls, mpls)
			if err != nil {
				return nil, err
			}
			obj.AddressType = addressTypePrefixList.String()
			obj.CidrBlock = aws.ToString(prefixList.PrefixListId) + "/" + name
			info = append(info, obj)
		}
	}
	return info, nil
}

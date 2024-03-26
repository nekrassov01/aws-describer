package describer

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func (a *app) doEc2Instance(c *cli.Context) error {
	switch a.dest.join {
	case ec2InstanceActionMemberDefault.String(), "":
		return a.doInstanceInfo(c)
	case ec2InstanceActionMemberSecurityGroup.String():
		return a.doInstanceSecurityGroupInfo(c)
	case ec2InstanceActionMemberRoute.String():
		return a.doInstanceRouteInfo(c)
	case ec2InstanceActionMemberStorage.String():
		return a.doInstanceStorageInfo(c)
	case ec2InstanceActionMemberBackup.String():
		return a.doInstanceBackupInfo(c)
	case ec2InstanceActionMemberLoadBalancer.String():
		return a.doInstanceLoadBalancerInfo(c)
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", a.dest.join, strings.Join(ec2InstanceActionMembers, "|"))
	}
}

func (a *app) doEc2Image(c *cli.Context) error {
	switch a.dest.join {
	case imageInfoActionMemberDefault.String(), "":
		return a.doImageInfo(c)
	case imageInfoActionMemberBackup.String():
		return a.doImageBackupInfo(c)
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", a.dest.join, strings.Join(ec2ImageActionMembers, "|"))
	}
}

func (a *app) doEc2SecurityGroup(c *cli.Context) error {
	switch a.dest.join {
	case ec2SecurityGroupActionMemberDefault.String(), "":
		return a.doSecurityGroupInfo(c)
	case ec2SecurityGroupActionMemberPermissions.String():
		return a.doSecurityGroupPermissionsInfo(c)
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", a.dest.join, strings.Join(ec2SecurityGroupActionMembers, "|"))
	}
}

func (a *app) doEc2Vpc(c *cli.Context) error {
	switch a.dest.join {
	case ec2VpcActionMemberDefault.String(), "":
		return a.doVpcInfo(c)
	case ec2VpcActionMemberAttribute.String():
		return a.doVpcAttributeInfo(c)
	case ec2VpcActionMemberCidr.String():
		return a.doVpcCidrInfo(c)
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", a.dest.join, strings.Join(ec2VpcActionMembers, "|"))
	}
}

func (a *app) doEc2Subnet(c *cli.Context) error {
	switch a.dest.join {
	case ec2SubnetActionMemberDefault.String(), "":
		return a.doSubnetInfo(c)
	case ec2SubnetActionMemberRoute.String():
		return a.doSubnetRouteInfo(c)
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", a.dest.join, strings.Join(ec2SubnetActionMembers, "|"))
	}
}

func (a *app) doEc2RouteTable(c *cli.Context) error {
	switch a.dest.join {
	case ec2RouteTableActionMemberDefault.String(), "":
		return a.doRouteTableInfo(c)
	case ec2RouteTableActionMemberPermissions.String():
		return a.doRouteTableAssociationInfo(c)
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", a.dest.join, strings.Join(ec2RouteTableActionMembers, "|"))
	}
}

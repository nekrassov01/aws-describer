package describer

type ec2InstanceActionMember int

const (
	ec2InstanceActionMemberDefault ec2InstanceActionMember = iota
	ec2InstanceActionMemberSecurityGroup
	ec2InstanceActionMemberRoute
	ec2InstanceActionMemberStorage
	ec2InstanceActionMemberBackup
	ec2InstanceActionMemberLoadBalancer
)

var ec2InstanceActionMembers = []string{
	"default",
	"sg",
	"route",
	"storage",
	"backup",
	"lb",
}

func (m ec2InstanceActionMember) String() string {
	if m >= 0 && int(m) < len(ec2InstanceActionMembers) {
		return ec2InstanceActionMembers[m]
	}
	return ""
}

type ec2ImageActionMember int

const (
	imageInfoActionMemberDefault ec2ImageActionMember = iota
	imageInfoActionMemberBackup
)

var ec2ImageActionMembers = []string{
	"default",
	"backup",
}

func (m ec2ImageActionMember) String() string {
	if m >= 0 && int(m) < len(ec2ImageActionMembers) {
		return ec2ImageActionMembers[m]
	}
	return ""
}

type ec2SecurityGroupActionMember int

const (
	ec2SecurityGroupActionMemberDefault ec2SecurityGroupActionMember = iota
	ec2SecurityGroupActionMemberPermissions
)

var ec2SecurityGroupActionMembers = []string{
	"default",
	"perms",
}

func (m ec2SecurityGroupActionMember) String() string {
	if m >= 0 && int(m) < len(ec2SecurityGroupActionMembers) {
		return ec2SecurityGroupActionMembers[m]
	}
	return ""
}

type ec2VpcActionMember int

const (
	ec2VpcActionMemberDefault ec2VpcActionMember = iota
	ec2VpcActionMemberAttribute
	ec2VpcActionMemberCidr
)

var ec2VpcActionMembers = []string{
	"default",
	"attr",
	"cidr",
}

func (m ec2VpcActionMember) String() string {
	if m >= 0 && int(m) < len(ec2VpcActionMembers) {
		return ec2VpcActionMembers[m]
	}
	return ""
}

type ec2SubnetActionMember int

const (
	ec2SubnetActionMemberDefault ec2SubnetActionMember = iota
	ec2SubnetActionMemberRoute
)

var ec2SubnetActionMembers = []string{
	"default",
	"route",
}

func (m ec2SubnetActionMember) String() string {
	if m >= 0 && int(m) < len(ec2SubnetActionMembers) {
		return ec2SubnetActionMembers[m]
	}
	return ""
}

type ec2RouteTableActionMember int

const (
	ec2RouteTableActionMemberDefault ec2RouteTableActionMember = iota
	ec2RouteTableActionMemberPermissions
)

var ec2RouteTableActionMembers = []string{
	"default",
	"assoc",
}

func (m ec2RouteTableActionMember) String() string {
	if m >= 0 && int(m) < len(ec2RouteTableActionMembers) {
		return ec2RouteTableActionMembers[m]
	}
	return ""
}

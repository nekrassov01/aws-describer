package ec2

type addressType int

const (
	addressTypeIpv4 addressType = iota
	addressTypeIpv6
	addressTypeSecurityGroup
	addressTypePrefixList
)

var addressTypes = []string{
	"Ipv4",
	"Ipv6",
	"SecurityGroup",
	"PrefixList",
}

func (a addressType) String() string {
	if a >= 0 && int(a) < len(addressTypes) {
		return addressTypes[a]
	}
	return ""
}

type targetType int

const (
	targetTypeLocal targetType = iota
	targetTypeInternetGateway
	targetTypeVpnGateway
	targetTypeNatGateway
	targetTypeVpcEndpoint
	targetTypeVpcPeeringConnection
	targetTypeTransitGateway
	targetTypeEgressOnlyInternetGateway
	targetTypeCarrierGateway
	targetTypeInstance
	targetTypeNetworkInterface
	targetTypeOutpostLocalGateway
	targetTypeCoreNetwork
	targetTypeOther
)

var targetTypes = []string{
	"Local",
	"InternetGateway",
	"VpnGateway",
	"NatGateway",
	"VpcEndpoint",
	"VpcPeeringConnection",
	"TransitGateway",
	"EgressOnlyInternetGateway",
	"CarrierGateway",
	"Instance",
	"NetworkInterface",
	"OutpostLocalGateway",
	"CoreNetwork",
	"Other",
}

func (t targetType) String() string {
	if t >= 0 && int(t) < len(targetTypes) {
		return targetTypes[t]
	}
	return ""
}

type flowDirection int

const (
	Ingress flowDirection = iota
	Egress
)

var flowDirections = []string{
	"Ingress",
	"Egress",
}

func (f flowDirection) String() string {
	if f >= 0 && int(f) < len(flowDirections) {
		return flowDirections[f]
	}
	return ""
}

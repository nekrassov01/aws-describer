package ec2

import (
	"context"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/nekrassov01/aws-describer/internal/api/elb"
	"github.com/nekrassov01/aws-describer/internal/api/elbv2"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func CreateDescribeInstancesInput(ids, names []string, filters []types.Filter, defaultFilter bool) *ec2.DescribeInstancesInput {
	var f []types.Filter
	if len(ids) > 0 {
		f = append(f, types.Filter{
			Name:   aws.String("instance-id"),
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
			Name:   aws.String("instance-state-name"),
			Values: instanceValidState,
		})
	}
	return &ec2.DescribeInstancesInput{
		Filters: f,
	}
}

type InstanceInfo struct {
	InstanceId       string
	InstanceName     string
	InstanceType     types.InstanceType
	PrivateIpAddress string
	PublicIpAddress  string
	Platform         string
	State            types.InstanceStateName
	AvailabilityZone string
}

func GetInstanceInfo(ich chan<- InstanceInfo, reservations []types.Reservation) {
	for _, r := range reservations {
		for _, i := range r.Instances {
			ich <- InstanceInfo{
				InstanceId:       aws.ToString(i.InstanceId),
				InstanceName:     getEc2NameTagValue(i.Tags),
				InstanceType:     i.InstanceType,
				PrivateIpAddress: aws.ToString(i.PrivateIpAddress),
				PublicIpAddress:  aws.ToString(i.PublicIpAddress),
				AvailabilityZone: aws.ToString(i.Placement.AvailabilityZone),
				Platform:         aws.ToString(i.PlatformDetails),
				State:            i.State.Name,
			}
		}
	}
}

type InstanceSecurityGroupInfo struct {
	InstanceId        string
	InstanceName      string
	VpcId             string
	VpcName           string
	SecurityGroupId   string
	SecurityGroupName string
	FlowDirection     string
	IpProtocol        string
	FromPort          int32
	ToPort            int32
	AddressType       string
	CidrBlock         string
	AvailabilityZone  string
}

func FetchDataForInstanceSecurityGroupInfo(ctx context.Context, l *rate.Limiter, client IEc2Client, region string) (map[string]types.SecurityGroup, map[string]types.Vpc, map[string]types.PrefixList, map[string]types.ManagedPrefixList, error) {
	segs := make(map[string]types.SecurityGroup)
	vpcs := make(map[string]types.Vpc)
	upls := make(map[string]types.PrefixList)
	mpls := make(map[string]types.ManagedPrefixList)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		segs, err = client.FetchSecurityGroups(ctx, region)
		return err
	})
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
		return nil, nil, nil, nil, err
	}
	return segs, vpcs, upls, mpls, nil
}

func GetInstanceSecurityGroupInfo(ich chan<- InstanceSecurityGroupInfo, reservations []types.Reservation, region string, segs map[string]types.SecurityGroup, vpcs map[string]types.Vpc, upls map[string]types.PrefixList, mpls map[string]types.ManagedPrefixList) error {
	for _, r := range reservations {
		for _, i := range r.Instances {
			name := getEc2NameTagValue(i.Tags)
			vpcId := aws.ToString(i.VpcId)
			vpc, err := findEc2VpcById(vpcId, vpcs)
			if err != nil {
				return err
			}
			vpcName := getEc2NameTagValue(vpc.Tags)
			for _, seg := range i.SecurityGroups {
				obj := InstanceSecurityGroupInfo{
					InstanceId:        aws.ToString(i.InstanceId),
					InstanceName:      name,
					AvailabilityZone:  aws.ToString(i.Placement.AvailabilityZone),
					VpcId:             vpcId,
					VpcName:           vpcName,
					SecurityGroupId:   aws.ToString(seg.GroupId),
					SecurityGroupName: aws.ToString(seg.GroupName),
				}
				sg, err := findEc2SecurityGroupById(aws.ToString(seg.GroupId), segs)
				if err != nil {
					return err
				}
				if err = appendToInstanceSecurityGroupInfo(ich, obj, *sg, vpcs, upls, mpls, false, region); err != nil {
					return err
				}
				if err = appendToInstanceSecurityGroupInfo(ich, obj, *sg, vpcs, upls, mpls, true, region); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func appendToInstanceSecurityGroupInfo(ich chan<- InstanceSecurityGroupInfo, obj InstanceSecurityGroupInfo, sg types.SecurityGroup, vpcs map[string]types.Vpc, upls map[string]types.PrefixList, mpls map[string]types.ManagedPrefixList, isEgress bool, region string) error {
	perms, err := handleSecurityGroupPermissionsInfo(sg, vpcs, upls, mpls, isEgress, region)
	if err != nil {
		return err
	}
	for _, perm := range perms {
		obj.FlowDirection = perm.FlowDirection
		obj.IpProtocol = perm.IpProtocol
		obj.FromPort = perm.FromPort
		obj.ToPort = perm.ToPort
		obj.AddressType = perm.AddressType
		obj.CidrBlock = perm.CidrBlock
		ich <- obj
	}
	return nil
}

type InstanceRouteInfo struct {
	InstanceId       string
	InstanceName     string
	VpcId            string
	VpcName          string
	SubnetId         string
	SubnetName       string
	AvailabilityZone string
	RouteTableId     string
	RouteTableName   string
	DestinationType  string
	Destination      string
	TargetType       string
	Target           string
	Region           string
}

func FetchDataForInstanceRouteInfo(ctx context.Context, l *rate.Limiter, client IEc2Client, region string) (map[string]types.Vpc, map[string]types.Subnet, map[string]types.RouteTable, error) {
	vpcs := make(map[string]types.Vpc)
	sbns := make(map[string]types.Subnet)
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
		sbns, err = client.FetchSubnets(ctx, region)
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
		return nil, nil, nil, err
	}
	return vpcs, sbns, rtbs, nil
}

func GetInstanceRouteInfo(ich chan<- InstanceRouteInfo, reservations []types.Reservation, region string, vpcs map[string]types.Vpc, sbns map[string]types.Subnet, rtbs map[string]types.RouteTable) error {
	for _, r := range reservations {
		for _, i := range r.Instances {
			vpc, err := findEc2VpcById(aws.ToString(i.VpcId), vpcs)
			if err != nil {
				return err
			}
			sbn, err := findEc2SubnetById(aws.ToString(i.SubnetId), sbns)
			if err != nil {
				return err
			}
			rtb, err := findEc2RouteTableBySubnet(*sbn, rtbs)
			if err != nil {
				return err
			}
			obj := InstanceRouteInfo{
				InstanceId:       aws.ToString(i.InstanceId),
				InstanceName:     getEc2NameTagValue(i.Tags),
				AvailabilityZone: aws.ToString(i.Placement.AvailabilityZone),
				VpcId:            aws.ToString(i.VpcId),
				VpcName:          getEc2NameTagValue(vpc.Tags),
				SubnetId:         aws.ToString(i.SubnetId),
				SubnetName:       getEc2NameTagValue(sbn.Tags),
				Region:           region,
			}
			routes, err := handleRoutes(*rtb, vpcs, region)
			if err != nil {
				return err
			}
			for _, route := range routes {
				obj.RouteTableId = route.RouteTableId
				obj.RouteTableName = route.RouteTableName
				obj.DestinationType = route.DestinationType
				obj.Destination = route.Destination
				obj.TargetType = route.TargetType
				obj.Target = route.Target
				ich <- obj
			}
		}
	}
	return nil
}

type InstanceStorageInfo struct {
	InstanceId          string
	InstanceName        string
	DeviceName          string
	DeleteOnTermination bool
	VolumeId            string
	VolumeName          string
	VolumeType          types.VolumeType
	VolumeSize          int32
	IOPS                int32
	Encrypted           bool
	AvailabilityZone    string
}

func GetInstanceStorageInfo(ich chan<- InstanceStorageInfo, reservations []types.Reservation, vols map[string]types.Volume) {
	for _, r := range reservations {
		for _, i := range r.Instances {
			for _, bdm := range i.BlockDeviceMappings {
				volId := aws.ToString(bdm.Ebs.VolumeId)
				vol := findEc2VolumeById(volId, vols)
				if vol != nil {
					ich <- InstanceStorageInfo{
						InstanceId:          aws.ToString(i.InstanceId),
						InstanceName:        getEc2NameTagValue(i.Tags),
						AvailabilityZone:    aws.ToString(i.Placement.AvailabilityZone),
						DeviceName:          aws.ToString(bdm.DeviceName),
						DeleteOnTermination: aws.ToBool(bdm.Ebs.DeleteOnTermination),
						VolumeId:            volId,
						VolumeName:          getEc2NameTagValue(vol.Tags),
						VolumeType:          vol.VolumeType,
						VolumeSize:          aws.ToInt32(vol.Size),
						IOPS:                aws.ToInt32(vol.Iops),
						Encrypted:           aws.ToBool(vol.Encrypted),
					}
				}
			}
		}
	}
}

type InstanceBackupInfo struct {
	InstanceId          string
	InstanceName        string
	ImageId             string
	ImageName           string
	ImageOwner          string
	DeleteOnTermination bool
	VolumeId            string
	VolumeName          string
	SnapshotId          string
	SnapshotName        string
	AvailabilityZone    string
}

func FetchDataForInstanceBackupInfo(ctx context.Context, l *rate.Limiter, client IEc2Client, region string) (map[string]types.Image, map[string]types.Snapshot, map[string]types.Volume, error) {
	var imgs map[string]types.Image
	var snps map[string]types.Snapshot
	var vols map[string]types.Volume
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		imgs, err = client.FetchImages(ctx, region)
		return err
	})
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		snps, err = client.FetchSnapshots(ctx, region)
		return err
	})
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		vols, err = client.FetchVolumes(ctx, region)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, nil, nil, err
	}
	return imgs, snps, vols, nil
}

func GetInstanceBackupInfo(ich chan InstanceBackupInfo, reservations []types.Reservation, imgs map[string]types.Image, snps map[string]types.Snapshot, vols map[string]types.Volume) {
	for _, r := range reservations {
		for _, i := range r.Instances {
			var imageId, imageName, imageOwner string
			imgId := aws.ToString(i.ImageId)
			img := findEc2ImageById(imgId, imgs)
			if img != nil {
				imageId = imgId
				imageName = aws.ToString(img.Name)
				imageOwner = aws.ToString(img.OwnerId)
			}
			for _, bdm := range i.BlockDeviceMappings {
				var volumeId, volumeName, snapshotId, snapshotName string
				if bdm.Ebs != nil && bdm.Ebs.VolumeId != nil {
					volId := aws.ToString(bdm.Ebs.VolumeId)
					vol := findEc2VolumeById(volId, vols)
					if vol != nil {
						volumeId = volId
						volumeName = getEc2NameTagValue(vol.Tags)
						snpId := aws.ToString(vol.SnapshotId)
						snp := findEc2SnapshotById(snpId, snps)
						if snp != nil {
							snapshotId = snpId
							snapshotName = getEc2NameTagValue(snp.Tags)
						}
					}
					ich <- InstanceBackupInfo{
						InstanceId:          aws.ToString(i.InstanceId),
						InstanceName:        getEc2NameTagValue(i.Tags),
						AvailabilityZone:    aws.ToString(i.Placement.AvailabilityZone),
						ImageId:             imageId,
						ImageName:           imageName,
						ImageOwner:          imageOwner,
						DeleteOnTermination: aws.ToBool(bdm.Ebs.DeleteOnTermination),
						SnapshotId:          snapshotId,
						SnapshotName:        snapshotName,
						VolumeId:            volumeId,
						VolumeName:          volumeName,
					}
				}
			}
		}
	}
}

type InstanceLoadBalancerInfo struct {
	InstanceId       string
	InstanceName     string
	AvailabilityZone string
	AttachedLB       []string
	AttachedTG       []string
}

func FetchDataForInstanceLoadBalancerInfo(ctx context.Context, l *rate.Limiter, cfg *aws.Config, region string, ids, names []string) ([]string, map[string][]string, map[string][]string, error) {
	var sv1, sv2 []string
	var mv1, mv2 map[string][]string
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		client := elb.NewElbClient(cfg)
		sv1, mv1, err = client.FetchTargets(ctx, region, true)
		return err
	})
	eg.Go(func() error {
		if err := l.Wait(ctx); err != nil {
			return err
		}
		var err error
		client := elbv2.NewElbClient(cfg)
		sv2, mv2, err = client.FetchTargets(ctx, region, "instance", true)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, nil, nil, err
	}
	if (len(sv1) > 0 || len(sv2) > 0) && (len(ids) == 0 && len(names) == 0) {
		sv1 := append(sv1, sv2...)
		slices.Sort(sv1)
		return slices.Compact(sv1), mv1, mv2, nil
	}
	return ids, mv1, mv2, nil
}

func GetInstanceLoadBalancerInfo(ich chan<- InstanceLoadBalancerInfo, reservations []types.Reservation, idmv1, idmv2 map[string][]string) error {
	for _, r := range reservations {
		for _, i := range r.Instances {
			id := aws.ToString(i.InstanceId)
			lbs, ok := idmv1[id]
			if !ok {
				lbs = nil
			}
			tgs, ok := idmv2[id]
			if !ok {
				tgs = nil
			}
			if len(lbs) == 0 && len(tgs) == 0 {
				return nil
			}
			ich <- InstanceLoadBalancerInfo{
				InstanceId:       aws.ToString(i.InstanceId),
				InstanceName:     getEc2NameTagValue(i.Tags),
				AvailabilityZone: aws.ToString(i.Placement.AvailabilityZone),
				AttachedLB:       lbs,
				AttachedTG:       tgs,
			}
		}
	}
	return nil
}

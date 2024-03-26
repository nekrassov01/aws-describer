package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func CreateDescribeImagesInput(ids, names []string, filters []types.Filter, defaultFilter bool) *ec2.DescribeImagesInput {
	var f []types.Filter
	if len(ids) > 0 {
		f = append(f, types.Filter{
			Name:   aws.String("image-id"),
			Values: ids,
		})
	}
	if len(names) > 0 {
		f = append(f, types.Filter{
			Name:   aws.String("name"),
			Values: names,
		})
	}
	if len(filters) > 0 {
		f = append(f, filters...)
	}
	input := &ec2.DescribeImagesInput{
		Filters: f,
	}
	if defaultFilter {
		input.Owners = []string{"self"}
	}
	return input
}

type ImageInfo struct {
	ImageId      string
	ImageName    string
	ImageOwner   string
	CreationDate string
	Architecture types.ArchitectureType
	PlatForm     string
	EnaSupport   bool
	Public       bool
	State        types.State
	Region       string
}

func GetImageInfo(ich chan<- ImageInfo, images []types.Image, region string) {
	for _, image := range images {
		ich <- ImageInfo{
			ImageId:      aws.ToString(image.ImageId),
			ImageName:    aws.ToString(image.Name),
			ImageOwner:   aws.ToString(image.OwnerId),
			CreationDate: aws.ToString(image.CreationDate),
			Architecture: types.ArchitectureType(image.Architecture),
			PlatForm:     aws.ToString(image.PlatformDetails),
			EnaSupport:   aws.ToBool(image.EnaSupport),
			Public:       aws.ToBool(image.Public),
			State:        types.State(image.State),
			Region:       region,
		}
	}
}

type ImageBackupInfo struct {
	ImageId             string
	ImageName           string
	ImageOwner          string
	DeleteOnTermination bool
	SnapshotId          string
	SnapshotName        string
	VolumeId            string
	VolumeName          string
	Region              string
}

func FetchDataForImageBackupInfo(ctx context.Context, l *rate.Limiter, client IEc2Client, region string) (map[string]types.Snapshot, map[string]types.Volume, error) {
	var snps map[string]types.Snapshot
	var vols map[string]types.Volume
	eg, ctx := errgroup.WithContext(ctx)
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
		return nil, nil, err
	}
	return snps, vols, nil
}

func GetImageBackupInfo(ich chan<- ImageBackupInfo, images []types.Image, region string, snps map[string]types.Snapshot, vols map[string]types.Volume) {
	for _, image := range images {
		var bds []*types.EbsBlockDevice
		for _, bdm := range image.BlockDeviceMappings {
			if bdm.Ebs != nil {
				bds = append(bds, bdm.Ebs)
			}
		}
		for _, bd := range bds {
			var snapshotId, snapshotName, volumeId, volumeName string
			if bd != nil && bd.SnapshotId != nil {
				snpId := aws.ToString(bd.SnapshotId)
				snp := findEc2SnapshotById(snpId, snps)
				if snp != nil {
					snapshotId = snpId
					snapshotName = getEc2NameTagValue(snp.Tags)
					volId := aws.ToString(snp.VolumeId)
					vol := findEc2VolumeById(aws.ToString(snp.VolumeId), vols)
					if vol != nil {
						volumeId = volId
						volumeName = getEc2NameTagValue(vol.Tags)
					}
				}
				ich <- ImageBackupInfo{
					ImageId:             aws.ToString(image.ImageId),
					ImageName:           aws.ToString(image.Name),
					ImageOwner:          aws.ToString(image.OwnerId),
					DeleteOnTermination: aws.ToBool(bd.DeleteOnTermination),
					SnapshotId:          snapshotId,
					SnapshotName:        snapshotName,
					VolumeId:            volumeId,
					VolumeName:          volumeName,
					Region:              region,
				}
			}
		}
	}
}

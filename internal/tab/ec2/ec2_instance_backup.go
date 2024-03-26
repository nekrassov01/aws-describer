// Code generated by table/gen.go. DO NOT EDIT.

package ec2

import (
	"sort"

	"github.com/nekrassov01/aws-describer/internal/api/ec2"
	"github.com/nekrassov01/aws-describer/internal/tab"
)

func PrintInstanceBackupInfo(info []ec2.InstanceBackupInfo, output string, header bool, mergeFields, ignoreFields []int) error {
	sort.SliceStable(info, func(i, j int) bool {
		if info[i].AvailabilityZone != info[j].AvailabilityZone {
			return info[i].AvailabilityZone < info[j].AvailabilityZone
		}
		if info[i].InstanceName != info[j].InstanceName {
			return info[i].InstanceName < info[j].InstanceName
		}
		if info[i].ImageOwner != info[j].ImageOwner {
			return info[i].ImageOwner < info[j].ImageOwner
		}
		if info[i].ImageName != info[j].ImageName {
			return info[i].ImageName < info[j].ImageName
		}
		if info[i].SnapshotName != info[j].SnapshotName {
			return info[i].SnapshotName < info[j].SnapshotName
		}
		return info[i].VolumeName < info[j].VolumeName
	})
	if len(mergeFields) == 0 {
		mergeFields = []int{0, 1, 2, 3, 4}
	}
	if err := tab.PrintTable(info, output, header, mergeFields, ignoreFields); err != nil {
		return err
	}
	return nil
}

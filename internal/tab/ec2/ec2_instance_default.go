// Code generated by table/gen.go. DO NOT EDIT.

package ec2

import (
	"sort"

	"github.com/nekrassov01/aws-describer/internal/api/ec2"
	"github.com/nekrassov01/aws-describer/internal/tab"
)

func PrintInstanceInfo(info []ec2.InstanceInfo, output string, header bool, mergeFields, ignoreFields []int) error {
	sort.SliceStable(info, func(i, j int) bool {
		if info[i].AvailabilityZone != info[j].AvailabilityZone {
			return info[i].AvailabilityZone < info[j].AvailabilityZone
		}
		if info[i].InstanceName != info[j].InstanceName {
			return info[i].InstanceName < info[j].InstanceName
		}
		if info[i].InstanceType != info[j].InstanceType {
			return info[i].InstanceType < info[j].InstanceType
		}
		if info[i].PrivateIpAddress != info[j].PrivateIpAddress {
			return info[i].PrivateIpAddress < info[j].PrivateIpAddress
		}
		return info[i].PrivateIpAddress < info[j].PrivateIpAddress
	})
	if err := tab.PrintTable(info, output, header, mergeFields, ignoreFields); err != nil {
		return err
	}
	return nil
}
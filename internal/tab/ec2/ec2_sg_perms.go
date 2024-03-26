// Code generated by table/gen.go. DO NOT EDIT.

package ec2

import (
	"sort"

	"github.com/nekrassov01/aws-describer/internal/api/ec2"
	"github.com/nekrassov01/aws-describer/internal/tab"
)

func PrintSecurityGroupPermissionsInfo(info []ec2.SecurityGroupPermissionsInfo, output string, header bool, mergeFields, ignoreFields []int) error {
	sort.SliceStable(info, func(i, j int) bool {
		if info[i].Region != info[j].Region {
			return info[i].Region < info[j].Region
		}
		if info[i].SecurityGroupName != info[j].SecurityGroupName {
			return info[i].SecurityGroupName < info[j].SecurityGroupName
		}
		if info[i].VpcName != info[j].VpcName {
			return info[i].VpcName < info[j].VpcName
		}
		if info[i].FlowDirection != info[j].FlowDirection {
			return info[i].FlowDirection > info[j].FlowDirection
		}
		if info[i].IpProtocol != info[j].IpProtocol {
			return info[i].IpProtocol < info[j].IpProtocol
		}
		if info[i].FromPort != info[j].FromPort {
			return info[i].FromPort < info[j].FromPort
		}
		if info[i].ToPort != info[j].ToPort {
			return info[i].ToPort < info[j].ToPort
		}
		return info[i].AddressType < info[j].AddressType
	})
	if len(mergeFields) == 0 {
		mergeFields = []int{0, 1, 2, 3}
	}
	if err := tab.PrintTable(info, output, header, mergeFields, ignoreFields); err != nil {
		return err
	}
	return nil
}

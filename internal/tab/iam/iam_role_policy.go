// Code generated by table/gen.go. DO NOT EDIT.

package iam

import (
	"sort"

	"github.com/nekrassov01/aws-describer/internal/api/iam"
	"github.com/nekrassov01/aws-describer/internal/tab"
)

func PrintRolePolicyInfo(info []iam.RolePolicyInfo, output string, header bool, mergeFields, ignoreFields []int, document bool) error {
	sort.SliceStable(info, func(i, j int) bool {
		if info[i].RoleName != info[j].RoleName {
			return info[i].RoleName < info[j].RoleName
		}
		return info[i].PolicyType < info[j].PolicyType
	})
	if len(mergeFields) == 0 {
		mergeFields = []int{0, 1}
	}
	if !document {
		ignoreFields = []int{5}
	}
	if err := tab.PrintTable(info, output, header, mergeFields, ignoreFields); err != nil {
		return err
	}
	return nil
}
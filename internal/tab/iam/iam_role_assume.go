// Code generated by table/gen.go. DO NOT EDIT.

package iam

import (
	"sort"

	"github.com/nekrassov01/aws-describer/internal/api/iam"
	"github.com/nekrassov01/aws-describer/internal/tab"
)

func PrintRoleAssumeInfo(info []iam.RoleAssumeInfo, output string, header bool, mergeFields, ignoreFields []int) error {
	sort.SliceStable(info, func(i, j int) bool {
		return info[i].RoleName < info[j].RoleName
	})
	if err := tab.PrintTable(info, output, header, mergeFields, ignoreFields); err != nil {
		return err
	}
	return nil
}

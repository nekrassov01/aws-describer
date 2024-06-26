// Code generated by table/gen.go. DO NOT EDIT.

package s3

import (
	"sort"

	"github.com/nekrassov01/aws-describer/internal/api/s3"
	"github.com/nekrassov01/aws-describer/internal/tab"
)

func PrintBucketInfo(info []s3.BucketInfo, output string, header bool, mergeFields, ignoreFields []int, document bool) error {
	sort.SliceStable(info, func(i, j int) bool {
		if info[i].BucketName != info[j].BucketName {
			return info[i].BucketName < info[j].BucketName
		}
		return info[i].Location < info[j].Location
	})
	if !document {
		ignoreFields = []int{3}
	}
	if err := tab.PrintTable(info, output, header, mergeFields, ignoreFields); err != nil {
		return err
	}
	return nil
}

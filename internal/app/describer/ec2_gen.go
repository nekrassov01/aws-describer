//go:build ignore

package main

import (
	"log"

	"github.com/nekrassov01/aws-describer/internal/tmpl"
)

type templateData struct {
	Name             string
	InputFuncName    string
	DescribeFuncName string
	PrintFuncName    string
}

func gen(filePath string, data templateData) error {
	template := `// Code generated by describer/ec2_gen.go. DO NOT EDIT.

package describer

import (
	ec2api "github.com/nekrassov01/aws-describer/internal/api/ec2"
	ec2tab "github.com/nekrassov01/aws-describer/internal/tab/ec2"
	"github.com/urfave/cli/v2"
)

func (a *app) {{ .Name }}(c *cli.Context) error {
	filters, err := ec2api.ParseEc2Filters(a.dest.ec2Filter)
	if err != nil {
		return err
	}
	info, err := ec2api.{{ .DescribeFuncName }}(c.Context, a.config, a.flag.regions.GetDestination(), a.flag.ids.GetDestination(), a.flag.names.GetDestination(), filters, a.dest.ec2DefaultFilter)
	if err != nil {
		return err
	}
	if err := ec2tab.{{ .PrintFuncName }}(info, a.dest.output, a.dest.header, a.flag.merge.GetDestination(), a.flag.ignore.GetDestination()); err != nil {
		return err
	}
	return nil
}`
	if err := tmpl.RenderTemplate("ec2", template, filePath, data); err != nil {
		return err
	}
	return nil
}

func main() {
	params := []struct {
		filePath string
		data     templateData
	}{
		{
			filePath: "ec2_instance_default.go",
			data: templateData{
				Name:             "doInstanceInfo",
				InputFuncName:    "CreateDescribeInstancesInput",
				DescribeFuncName: "DescribeInstanceInfo",
				PrintFuncName:    "PrintInstanceInfo",
			},
		},
		{
			filePath: "ec2_instance_sg.go",
			data: templateData{
				Name:             "doInstanceSecurityGroupInfo",
				InputFuncName:    "CreateDescribeInstancesInput",
				DescribeFuncName: "DescribeInstanceSecurityGroupInfo",
				PrintFuncName:    "PrintInstanceSecurityGroupInfo",
			},
		},
		{
			filePath: "ec2_instance_rtb.go",
			data: templateData{
				Name:             "doInstanceRouteInfo",
				InputFuncName:    "CreateDescribeInstancesInput",
				DescribeFuncName: "DescribeInstanceRouteInfo",
				PrintFuncName:    "PrintInstanceRouteInfo",
			},
		},
		{
			filePath: "ec2_instance_storage.go",
			data: templateData{
				Name:             "doInstanceStorageInfo",
				InputFuncName:    "CreateDescribeInstancesInput",
				DescribeFuncName: "DescribeInstanceStorageInfo",
				PrintFuncName:    "PrintInstanceStorageInfo",
			},
		},
		{
			filePath: "ec2_instance_backup.go",
			data: templateData{
				Name:             "doInstanceBackupInfo",
				InputFuncName:    "CreateDescribeInstancesInput",
				DescribeFuncName: "DescribeInstanceBackupInfo",
				PrintFuncName:    "PrintInstanceBackupInfo",
			},
		},
		{
			filePath: "ec2_instance_lb.go",
			data: templateData{
				Name:             "doInstanceLoadBalancerInfo",
				InputFuncName:    "CreateDescribeInstancesInput",
				DescribeFuncName: "DescribeInstanceLoadBalancerInfo",
				PrintFuncName:    "PrintInstanceLoadBalancerInfo",
			},
		},
		{
			filePath: "ec2_image_default.go",
			data: templateData{
				Name:             "doImageInfo",
				InputFuncName:    "CreateDescribeImagesInput",
				DescribeFuncName: "DescribeImageInfo",
				PrintFuncName:    "PrintImageInfo",
			},
		},
		{
			filePath: "ec2_image_backup.go",
			data: templateData{
				Name:             "doImageBackupInfo",
				InputFuncName:    "CreateDescribeImagesInput",
				DescribeFuncName: "DescribeImageBackupInfo",
				PrintFuncName:    "PrintImageBackupInfo",
			},
		},
		{
			filePath: "ec2_sg_default.go",
			data: templateData{
				Name:             "doSecurityGroupInfo",
				InputFuncName:    "CreateDescribeSecurityGroupsInput",
				DescribeFuncName: "DescribeSecurityGroupInfo",
				PrintFuncName:    "PrintSecurityGroupInfo",
			},
		},
		{
			filePath: "ec2_sg_perms.go",
			data: templateData{
				Name:             "doSecurityGroupPermissionsInfo",
				InputFuncName:    "CreateDescribeSecurityGroupsInput",
				DescribeFuncName: "DescribeSecurityGroupPermissionsInfo",
				PrintFuncName:    "PrintSecurityGroupPermissionsInfo",
			},
		},
		{
			filePath: "ec2_vpc_default.go",
			data: templateData{
				Name:             "doVpcInfo",
				InputFuncName:    "CreateDescribeVpcsInput",
				DescribeFuncName: "DescribeVpcInfo",
				PrintFuncName:    "PrintVpcInfo",
			},
		},
		{
			filePath: "ec2_vpc_attr.go",
			data: templateData{
				Name:             "doVpcAttributeInfo",
				InputFuncName:    "CreateDescribeVpcsInput",
				DescribeFuncName: "DescribeVpcAttributeInfo",
				PrintFuncName:    "PrintVpcAttributeInfo",
			},
		},
		{
			filePath: "ec2_vpc_cidr.go",
			data: templateData{
				Name:             "doVpcCidrInfo",
				InputFuncName:    "CreateDescribeVpcsInput",
				DescribeFuncName: "DescribeVpcCidrInfo",
				PrintFuncName:    "PrintVpcCidrInfo",
			},
		},
		{
			filePath: "ec2_subnet_default.go",
			data: templateData{
				Name:             "doSubnetInfo",
				InputFuncName:    "CreateDescribeSubnetsInput",
				DescribeFuncName: "DescribeSubnetInfo",
				PrintFuncName:    "PrintSubnetInfo",
			},
		},
		{
			filePath: "ec2_subnet_route.go",
			data: templateData{
				Name:             "doSubnetRouteInfo",
				InputFuncName:    "CreateDescribeSubnetsInput",
				DescribeFuncName: "DescribeSubnetRouteInfo",
				PrintFuncName:    "PrintSubnetRouteInfo",
			},
		},
		{
			filePath: "ec2_rtb_default.go",
			data: templateData{
				Name:             "doRouteTableInfo",
				InputFuncName:    "CreateDescribeRouteTablesInput",
				DescribeFuncName: "DescribeRouteTableInfo",
				PrintFuncName:    "PrintRouteTableInfo",
			},
		},
		{
			filePath: "ec2_rtb_assoc.go",
			data: templateData{
				Name:             "doRouteTableAssociationInfo",
				InputFuncName:    "CreateDescribeRouteTablesInput",
				DescribeFuncName: "DescribeRouteTableAssociationInfo",
				PrintFuncName:    "PrintRouteTableAssociationInfo",
			},
		},
	}
	for _, param := range params {
		if err := gen(param.filePath, param.data); err != nil {
			log.Fatal(err)
		}
	}
}
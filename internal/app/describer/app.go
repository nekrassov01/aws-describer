//go:generate go run ec2_gen.go
//go:generate go run iam_gen.go
//go:generate go run s3_gen.go
//go:generate gofmt -w .

package describer

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/nekrassov01/aws-describer/internal/api"
	"github.com/nekrassov01/aws-describer/internal/api/iam"
	"github.com/nekrassov01/mintab"
	"github.com/urfave/cli/v2"
)

const Name = "aws-describer"

type app struct {
	App    *cli.App
	config *aws.Config
	dest   dest
	flag   flag
}

type dest struct {
	join             string
	output           string
	profile          string
	region           string
	regions          cli.StringSlice
	ids              cli.StringSlice
	names            cli.StringSlice
	header           bool
	merge            cli.IntSlice
	ignore           cli.IntSlice
	document         bool
	documentFilter   cli.StringSlice
	ec2Filter        string
	ec2DefaultFilter bool
	iamPolicyScope   string
}

type flag struct {
	join             *cli.StringFlag
	output           *cli.StringFlag
	profile          *cli.StringFlag
	region           *cli.StringFlag
	regions          *cli.StringSliceFlag
	ids              *cli.StringSliceFlag
	names            *cli.StringSliceFlag
	header           *cli.BoolFlag
	merge            *cli.IntSliceFlag
	ignore           *cli.IntSliceFlag
	document         *cli.BoolFlag
	documentFilter   *cli.StringSliceFlag
	ec2Filter        *cli.StringFlag
	ec2DefaultFilter *cli.BoolFlag
	iamPolicyScope   *cli.StringFlag
}

func New() *app {
	a := app{}
	a.flag.output = &cli.StringFlag{
		Name:        "output",
		Aliases:     []string{"o"},
		Usage:       fmt.Sprintf("select output format: %s", strings.Join(mintab.Formats, "|")),
		Destination: &a.dest.output,
		Value:       mintab.FormatText.String(),
		EnvVars:     []string{strings.ToUpper(strings.ReplaceAll(Name, "-", "_")) + "_OUTPUT_FORMAT"},
	}
	a.flag.profile = &cli.StringFlag{
		Name:        "profile",
		Aliases:     []string{"p"},
		Usage:       "set aws profile",
		Destination: &a.dest.profile,
	}
	a.flag.region = &cli.StringFlag{
		Name:        "region",
		Aliases:     []string{"r"},
		Usage:       "set primary region",
		Destination: &a.dest.region,
		Value:       api.DefaultRegion,
		EnvVars:     []string{"AWS_DEFAULT_REGION"},
	}
	a.flag.regions = &cli.StringSliceFlag{
		Name:        "regions",
		Aliases:     []string{"R"},
		Usage:       "set target regions to request",
		Destination: &a.dest.regions,
		EnvVars:     []string{strings.ToUpper(strings.ReplaceAll(Name, "-", "_")) + "_TARGET_REGIONS"},
		DefaultText: "all enabled ec2 regions by default",
	}
	a.flag.regions.SetValue(api.DefaultTargetRegions)
	a.flag.ids = &cli.StringSliceFlag{
		Name:        "ids",
		Aliases:     []string{"i"},
		Usage:       "set resource ids",
		Destination: &a.dest.ids,
	}
	a.flag.names = &cli.StringSliceFlag{
		Name:        "names",
		Aliases:     []string{"n"},
		Usage:       "set resource names or name tags",
		Destination: &a.dest.names,
	}
	a.flag.header = &cli.BoolFlag{
		Name:        "header",
		Aliases:     []string{"H"},
		Usage:       "set whether to disable table header",
		Destination: &a.dest.header,
		Value:       true,
	}
	a.flag.merge = &cli.IntSliceFlag{
		Name:        "merge",
		Aliases:     []string{"M"},
		Usage:       "set column indexes to merge by value",
		Destination: &a.dest.merge,
	}
	a.flag.ignore = &cli.IntSliceFlag{
		Name:        "ignore",
		Aliases:     []string{"I"},
		Usage:       "set column indexes to exclude from output",
		Destination: &a.dest.ignore,
	}
	a.flag.document = &cli.BoolFlag{
		Name:        "document",
		Aliases:     []string{"d"},
		Usage:       "enable output of policy documents",
		Destination: &a.dest.document,
		Value:       false,
	}
	a.flag.documentFilter = &cli.StringSliceFlag{
		Name:        "document-filter",
		Aliases:     []string{"f"},
		Usage:       "set words to filter policy documents",
		Destination: &a.dest.documentFilter,
	}
	a.flag.ec2Filter = &cli.StringFlag{
		Name:        "filter",
		Aliases:     []string{"f"},
		Usage:       "set ec2 filter: '{name: \"key\", values: [\"value1\", \"value2\"]}'",
		Destination: &a.dest.ec2Filter,
	}
	a.flag.ec2DefaultFilter = &cli.BoolFlag{
		Name:        "default-filter",
		Aliases:     []string{"d"},
		Usage:       "set whether to disable default ec2 filter",
		Destination: &a.dest.ec2DefaultFilter,
		Value:       true,
	}
	a.flag.iamPolicyScope = &cli.StringFlag{
		Name:        "scope",
		Aliases:     []string{"s"},
		Usage:       fmt.Sprintf("select iam policy scope: %s", strings.Join(iam.PolicyScopeTypes, "|")),
		Destination: &a.dest.iamPolicyScope,
		Value:       iam.PolicyScopeTypeLocal.String(),
	}
	baseFlags := func(s []string) []cli.Flag {
		a.joinFlag(s)
		return []cli.Flag{
			a.flag.join,
			a.flag.output,
			a.flag.region,
			a.flag.profile,
			a.flag.header,
			a.flag.merge,
			a.flag.ignore,
		}
	}
	ec2DefaultFlags := func(s []string) []cli.Flag {
		return append(
			baseFlags(s),
			a.flag.regions,
			a.flag.ids,
			a.flag.names,
			a.flag.ec2Filter,
			a.flag.ec2DefaultFilter,
		)
	}
	iamDefaultFlags := func(s []string) []cli.Flag {
		return append(
			baseFlags(s),
			a.flag.ids,
			a.flag.names,
			a.flag.document,
			a.flag.documentFilter,
		)
	}
	iamScopeFlags := func(s []string) []cli.Flag {
		return append(
			iamDefaultFlags(s),
			a.flag.iamPolicyScope,
		)
	}
	s3DefaultFlags := func(s []string) []cli.Flag {
		return append(
			baseFlags(s),
			a.flag.regions,
			a.flag.names,
			a.flag.document,
			a.flag.documentFilter,
		)
	}
	a.App = &cli.App{
		Name:                 Name,
		Usage:                "AWS resources describer CLI",
		Version:              Version,
		Description:          "A cli application to join and list AWS resources with various other resources",
		HideHelpCommand:      true,
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:            "completion",
				Description:     "Generate completion scripts",
				Usage:           fmt.Sprintf("Generate completion scripts: %s", strings.Join(shells, "|")),
				UsageText:       fmt.Sprintf("%s completion shell", Name),
				HideHelpCommand: true,
				Action:          a.doCompletion,
			},
			{
				Name:            "ec2",
				Description:     "Invoke EC2 API and list resources in various output formats",
				Usage:           "Invoke EC2 API and list resources",
				UsageText:       fmt.Sprintf("%s ec2 command", Name),
				HideHelpCommand: true,
				Before:          a.doBefore,
				Subcommands: []*cli.Command{
					{
						Name:            "get-instances",
						Description:     "List EC2 instance info in combination with various resources",
						Usage:           "List EC2 instance info",
						UsageText:       fmt.Sprintf("%s ec2 get-instances", Name),
						HideHelpCommand: true,
						Flags:           ec2DefaultFlags(ec2InstanceActionMembers),
						Action:          a.doEc2Instance,
					},
					{
						Name:            "get-images",
						Description:     "List EC2 image info in combination with various resources",
						Usage:           "List EC2 image info",
						UsageText:       fmt.Sprintf("%s ec2 get-images", Name),
						HideHelpCommand: true,
						Flags:           ec2DefaultFlags(ec2ImageActionMembers),
						Action:          a.doEc2Image,
					},
					{
						Name:            "get-security-groups",
						Description:     "List EC2 security group info in combination with various resources",
						Usage:           "List EC2 security group info",
						UsageText:       fmt.Sprintf("%s ec2 get-security-groups", Name),
						HideHelpCommand: true,
						Flags:           ec2DefaultFlags(ec2SecurityGroupActionMembers),
						Action:          a.doEc2SecurityGroup,
					},
					{
						Name:            "get-vpcs",
						Description:     "List EC2 VPC info in combination with various resources",
						Usage:           "List EC2 VPC info",
						UsageText:       fmt.Sprintf("%s ec2 get-vpcs", Name),
						HideHelpCommand: true,
						Flags:           ec2DefaultFlags(ec2VpcActionMembers),
						Action:          a.doEc2Vpc,
					},
					{
						Name:            "get-subnets",
						Description:     "List EC2 subnet info in combination with various resources",
						Usage:           "List EC2 subnet info",
						UsageText:       fmt.Sprintf("%s ec2 get-subnets", Name),
						HideHelpCommand: true,
						Flags:           ec2DefaultFlags(ec2SubnetActionMembers),
						Action:          a.doEc2Subnet,
					},
					{
						Name:            "get-route-tables",
						Description:     "List EC2 route table info in combination with various resources",
						Usage:           "List EC2 route table info",
						UsageText:       fmt.Sprintf("%s ec2 get-route-tables", Name),
						HideHelpCommand: true,
						Flags:           ec2DefaultFlags(ec2RouteTableActionMembers),
						Action:          a.doEc2RouteTable,
					},
				},
			},
			{
				Name:            "iam",
				Description:     "Invoke IAM API and list resources in various output formats",
				Usage:           "Invoke IAM API and list resources",
				UsageText:       fmt.Sprintf("%s iam command", Name),
				HideHelpCommand: true,
				Before:          a.doBefore,
				Subcommands: []*cli.Command{
					{
						Name:            "get-users",
						Description:     "List IAM user info in combination with related info",
						Usage:           "List IAM user info",
						UsageText:       fmt.Sprintf("%s iam get-users", Name),
						HideHelpCommand: true,
						Flags:           iamDefaultFlags(iamUserActionMembers),
						Action:          a.doIamUser,
					},
					{
						Name:            "get-groups",
						Description:     "List IAM group info in combination with related info",
						Usage:           "List IAM group info",
						UsageText:       fmt.Sprintf("%s iam get-groups", Name),
						HideHelpCommand: true,
						Flags:           iamDefaultFlags(iamGroupActionMembers),
						Action:          a.doIamGroup,
					},
					{
						Name:            "get-roles",
						Description:     "List IAM role info in combination with related info",
						Usage:           "List IAM role info",
						UsageText:       fmt.Sprintf("%s iam get-roles", Name),
						HideHelpCommand: true,
						Flags:           iamDefaultFlags(iamRoleActionMembers),
						Action:          a.doIamRole,
					},
					{
						Name:            "get-policies",
						Description:     "List IAM policy info in combination with related info",
						Usage:           "List IAM policy info",
						UsageText:       fmt.Sprintf("%s iam get-policies", Name),
						HideHelpCommand: true,
						Flags:           iamScopeFlags(iamPolicyActionMembers),
						Action:          a.doIamPolicy,
					},
				},
			},
			{
				Name:            "s3",
				Description:     "Invoke S3 API and list resources in various output formats",
				Usage:           "Invoke S3 API and list resources",
				UsageText:       fmt.Sprintf("%s s3 command", Name),
				HideHelpCommand: true,
				Before:          a.doBefore,
				Subcommands: []*cli.Command{
					{
						Name:            "get-buckets",
						Description:     "List S3 bucket info in combination with related info",
						Usage:           "List S3 bucket info",
						UsageText:       fmt.Sprintf("%s s3 get-buckets", Name),
						HideHelpCommand: true,
						Flags:           s3DefaultFlags(s3BucketActionMembers),
						Action:          a.doS3Bucket,
					},
				},
			},
		},
	}
	return &a
}

func (a *app) joinFlag(s []string) {
	a.flag.join = &cli.StringFlag{
		Name:        "join",
		Aliases:     []string{"j"},
		Usage:       fmt.Sprintf("set info to be joined: %s", strings.Join(s, "|")),
		Destination: &a.dest.join,
	}
}

func (a *app) doBefore(c *cli.Context) error {
	cfg, err := api.LoadConfig(c.Context, a.dest.region, a.dest.profile)
	if err != nil {
		return err
	}
	a.config = cfg
	return nil
}

func (a *app) doCompletion(c *cli.Context) error {
	shell := c.Args().First()
	switch shell {
	case bash.String():
		fmt.Println(bashCompletion)
	case zsh.String():
		fmt.Println(zshCompletion)
	case pwsh.String():
		fmt.Println(pwshCompletion)
	default:
		return fmt.Errorf("%s: unsupported shell: valid values: %s", shell, strings.Join(shells, "|"))
	}
	return nil
}

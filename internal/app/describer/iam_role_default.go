// Code generated by describer/iam_gen.go. DO NOT EDIT.

package describer

import (
	"fmt"

	iamapi "github.com/nekrassov01/aws-describer/internal/api/iam"
	iamtab "github.com/nekrassov01/aws-describer/internal/tab/iam"
	"github.com/urfave/cli/v2"
)

func (a *app) doRoleInfo(c *cli.Context) error {
	if c.IsSet(a.flag.document.Name) || c.IsSet(a.flag.documentFilter.Name) {
		return fmt.Errorf("invalid args/flags combination: \"%s\" and \"%s\" are valid only when \"%s\" is selected at \"%s\"", a.flag.document.Name, a.flag.documentFilter.Name, iamGroupActionMemberPolicy.String(), a.flag.join.Name)
	}
	client := iamapi.NewIamClient(a.config)
	info, err := iamapi.ListRoleInfo(c.Context, client, a.flag.ids.GetDestination(), a.flag.names.GetDestination())
	if err != nil {
		return err
	}
	if err := iamtab.PrintRoleInfo(info, a.dest.output, a.dest.header, a.flag.merge.GetDestination(), a.flag.ignore.GetDestination()); err != nil {
		return err
	}
	return nil
}

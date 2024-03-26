// Code generated by describer/iam_gen.go. DO NOT EDIT.

package describer

import (
	"fmt"

	iamapi "github.com/nekrassov01/aws-describer/internal/api/iam"
	iamtab "github.com/nekrassov01/aws-describer/internal/tab/iam"
	"github.com/urfave/cli/v2"
)

func (a *app) doUserPolicyInfo(c *cli.Context) error {
	if !c.IsSet(a.flag.document.Name) && c.IsSet(a.flag.documentFilter.Name) {
		return fmt.Errorf("invalid args/flags combination: \"%s\" is valid only when \"%s\" is enabled", a.flag.documentFilter.Name, a.flag.document.Name)
	}
	client := iamapi.NewIamClient(a.config)
	info, err := iamapi.ListUserPolicyInfo(c.Context, client, a.flag.ids.GetDestination(), a.flag.names.GetDestination(), a.dest.document, a.flag.documentFilter.GetDestination())
	if err != nil {
		return err
	}
	if err := iamtab.PrintUserPolicyInfo(info, a.dest.output, a.dest.header, a.flag.merge.GetDestination(), a.flag.ignore.GetDestination(), a.dest.document); err != nil {
		return err
	}
	return nil
}

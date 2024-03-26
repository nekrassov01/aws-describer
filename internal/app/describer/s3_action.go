package describer

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func (a *app) doS3Bucket(c *cli.Context) error {
	switch a.dest.join {
	case s3BucketActionMemberDefault.String(), "":
		return a.doBucketInfo(c)
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", a.dest.join, strings.Join(s3BucketActionMembers, "|"))
	}
}

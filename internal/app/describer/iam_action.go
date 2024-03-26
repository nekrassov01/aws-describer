package describer

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func (a *app) doIamUser(c *cli.Context) error {
	switch a.dest.join {
	case iamUserActionMemberDefault.String(), "":
		return a.doUserInfo(c)
	case iamUserActionMemberPolicy.String():
		return a.doUserPolicyInfo(c)
	case iamUserActionMemberGroup.String():
		return a.doUserGroupInfo(c)
	case iamUserActionMemberAssociation.String():
		return a.doUserAssociationInfo(c)
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", a.dest.join, strings.Join(iamUserActionMembers, "|"))
	}
}

func (a *app) doIamGroup(c *cli.Context) error {
	switch a.dest.join {
	case iamGroupActionMemberDefault.String(), "":
		return a.doGroupInfo(c)
	case iamGroupActionMemberPolicy.String():
		return a.doGroupPolicyInfo(c)
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", a.dest.join, strings.Join(iamGroupActionMembers, "|"))
	}
}

func (a *app) doIamRole(c *cli.Context) error {
	switch a.dest.join {
	case iamRoleActionMemberDefault.String(), "":
		return a.doRoleInfo(c)
	case iamRoleActionMemberPolicy.String():
		return a.doRolePolicyInfo(c)
	case iamRoleActionMemberAssume.String():
		return a.doRoleAssumeInfo(c)
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", a.dest.join, strings.Join(iamRoleActionMembers, "|"))
	}
}

func (a *app) doIamPolicy(c *cli.Context) error {
	switch a.dest.join {
	case iamPolicyActionMemberDefault.String(), "":
		return a.doPolicyInfo(c)
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", a.dest.join, strings.Join(iamPolicyActionMembers, "|"))
	}
}

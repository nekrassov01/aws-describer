package describer

type iamUserActionMember int

const (
	iamUserActionMemberDefault iamUserActionMember = iota
	iamUserActionMemberPolicy
	iamUserActionMemberGroup
	iamUserActionMemberAssociation
)

var iamUserActionMembers = []string{
	"default",
	"policy",
	"group",
	"assoc",
}

func (m iamUserActionMember) String() string {
	if m >= 0 && int(m) < len(iamUserActionMembers) {
		return iamUserActionMembers[m]
	}
	return ""
}

type iamGroupActionMember int

const (
	iamGroupActionMemberDefault iamGroupActionMember = iota
	iamGroupActionMemberPolicy
)

var iamGroupActionMembers = []string{
	"default",
	"policy",
}

func (m iamGroupActionMember) String() string {
	if m >= 0 && int(m) < len(iamGroupActionMembers) {
		return iamGroupActionMembers[m]
	}
	return ""
}

type iamRoleActionMember int

const (
	iamRoleActionMemberDefault iamRoleActionMember = iota
	iamRoleActionMemberPolicy
	iamRoleActionMemberAssume
)

var iamRoleActionMembers = []string{
	"default",
	"policy",
	"assume",
}

func (m iamRoleActionMember) String() string {
	if m >= 0 && int(m) < len(iamRoleActionMembers) {
		return iamRoleActionMembers[m]
	}
	return ""
}

type iamPolicyActionMember int

const (
	iamPolicyActionMemberDefault iamPolicyActionMember = iota
)

var iamPolicyActionMembers = []string{
	"default",
}

func (m iamPolicyActionMember) String() string {
	if m >= 0 && int(m) < len(iamPolicyActionMembers) {
		return iamPolicyActionMembers[m]
	}
	return ""
}

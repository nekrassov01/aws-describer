package iam

type policyType int

const (
	policyTypeAttached policyType = iota
	policyTypeInline
)

var policyTypes = []string{
	"Attached",
	"Inline",
}

func (t policyType) String() string {
	if t >= 0 && int(t) < len(policyTypes) {
		return policyTypes[t]
	}
	return ""
}

type PolicyScopeType int

const (
	PolicyScopeTypeLocal PolicyScopeType = iota
	PolicyScopeTypeAws
)

var PolicyScopeTypes = []string{
	"local",
	"aws",
}

func (t PolicyScopeType) String() string {
	if t >= 0 && int(t) < len(PolicyScopeTypes) {
		return PolicyScopeTypes[t]
	}
	return ""
}

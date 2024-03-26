package describer

type s3BucketActionMember int

const (
	s3BucketActionMemberDefault s3BucketActionMember = iota
)

var s3BucketActionMembers = []string{
	"default",
}

func (m s3BucketActionMember) String() string {
	if m >= 0 && int(m) < len(s3BucketActionMembers) {
		return s3BucketActionMembers[m]
	}
	return ""
}

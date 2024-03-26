package describer

type shell int

const (
	bash shell = iota
	zsh
	pwsh
)

var shells = []string{
	"bash",
	"zsh",
	"pwsh",
}

func (s shell) String() string {
	if s >= 0 && int(s) < len(shells) {
		return shells[s]
	}
	return ""
}

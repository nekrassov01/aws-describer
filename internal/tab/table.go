//go:generate go run gen.go
//go:generate gofmt -w ec2/
//go:generate gofmt -w iam/
//go:generate gofmt -w s3/

package tab

import (
	"fmt"
	"os"
	"strings"

	"github.com/nekrassov01/mintab"
)

func PrintTable(info any, output string, header bool, mergeFields, ignoreFields []int) error {
	var o mintab.Format
	switch output {
	case mintab.FormatText.String():
		o = mintab.FormatText
	case mintab.FormatCompressedText.String():
		o = mintab.FormatCompressedText
	case mintab.FormatMarkdown.String():
		o = mintab.FormatMarkdown
	case mintab.FormatBacklog.String():
		o = mintab.FormatBacklog
	default:
		return fmt.Errorf("invalid value: %s: valid values: %s", output, strings.Join(mintab.Formats, "|"))
	}
	table := mintab.New(
		os.Stdout,
		mintab.WithFormat(o),
		mintab.WithHeader(header),
		mintab.WithMergeFields(mergeFields),
		mintab.WithIgnoreFields(ignoreFields),
	)
	if err := table.Load(info); err != nil {
		return fmt.Errorf("cannot output result: %w", err)
	}
	table.Out()
	return nil
}

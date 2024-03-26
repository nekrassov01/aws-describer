package ec2

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/google/go-jsonnet"
)

func ParseEc2Filters(s string) ([]types.Filter, error) {
	s, err := toJSON(toBracket(s))
	if err != nil {
		return nil, err
	}
	var filters []types.Filter
	if err := json.Unmarshal([]byte(s), &filters); err != nil {
		return nil, fmt.Errorf("cannot unmarshal value passed in filter: %w", err)
	}
	if err := validateFilters(filters); err != nil {
		return nil, fmt.Errorf("cannot parse value passed in filter: %w", err)
	}
	return filters, nil
}

func toBracket(s string) string {
	if !(strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")) {
		s = "[" + s + "]"
	}
	return s
}

func toJSON(s string) (string, error) {
	vm := jsonnet.MakeVM()
	j, err := vm.EvaluateAnonymousSnippet("", s)
	if err != nil {
		return "", fmt.Errorf("cannot convert to json from jsonnet code: %w", err)
	}
	return j, nil
}

func validateFilters(filters []types.Filter) error {
	for _, filter := range filters {
		if aws.ToString(filter.Name) == "" {
			return fmt.Errorf("empty [Nn]ame in filter string")
		}
		if len(filter.Values) == 0 {
			return fmt.Errorf("empty [Vv]alues in filter string")
		}
	}
	return nil
}

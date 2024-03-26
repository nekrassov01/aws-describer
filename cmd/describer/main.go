package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nekrassov01/aws-describer/internal/app/describer"
)

func main() {
	if err := describer.New().App.RunContext(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("%s: %w", describer.Name, err))
		os.Exit(1)
	}
}

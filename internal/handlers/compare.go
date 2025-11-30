package handlers

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func Compare(ctx context.Context, cmd *cli.Command) (map[string][]string, error) {
	fmt.Println("compare folders")
	return nil, nil
}

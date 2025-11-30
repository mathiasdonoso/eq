package handlers

import (
	"bufio"
	"context"
	"os"

	"github.com/mathiasdonoso/eq/pkg/printer"
	"github.com/urfave/cli/v3"
)

func Handler(ctx context.Context, cmd *cli.Command) error {
	var r map[string][]string
	var err error

	if cmd.Bool("files") {
		r, err = CompareFiles(ctx, cmd)
	} else {
		r, err = Compare(ctx, cmd)
	}

	if err != nil {
		return err
	}

	output := bufio.NewWriter(os.Stdout)
	printer.Print(output, r)
	defer output.Flush()

	return nil
}

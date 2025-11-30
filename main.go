package main

import (
	"context"
	"log"
	"os"

	"github.com/mathiasdonoso/eq/internal/handlers"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:        "eq",
		Description: "A fast CLI utility that detects duplicate files",
		Version:     "0.1.0",
		Authors:     []any{"Mathias Donoso <mathiasd88@gmail.com>"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "files",
				Value: false,
				Usage: "Compare explicit files only",
			},
			&cli.StringFlag{
				Name:  "hash",
				Value: "blake3",
				Usage: "Hash algorithm to use.\n\tSupported: sha256, blake3, xxh64",
			},
			&cli.BoolFlag{
				Name:  "verify",
				Value: false,
				Usage: "Perform byte-to-byte verification on\n\tfiles that share the same hash",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return handlers.Handler(ctx, cmd)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

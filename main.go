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
		Name:      "eq",
		Usage:     "A fast CLI utility that detects duplicate files",
		UsageText: "eq [paths...] [options]",
		Description: `eq scans files and directories, hashes their contents, and reports files that are
exact duplicates. It works recursively, supports multiple hashing algorithms, and
can optionally verify matches byte-by-byte. You can mix files and directories, and
eq will walk each path, skip symlinks, and group files by hash.`,
		Version: "0.1.0",
		Authors: []any{"Mathias Donoso <mathiasd88@gmail.com>"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "hash",
				Value: "blake3",
				Usage: "Hash algorithm to use.\n\tSupported: sha256, blake3, xxh64",
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

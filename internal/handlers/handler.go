package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mathiasdonoso/eq/pkg/hash"
	"github.com/mathiasdonoso/eq/pkg/printer"
	"github.com/urfave/cli/v3"
)

func Handler(ctx context.Context, cmd *cli.Command) error {
	r, err := Run(ctx, cmd)
	if err != nil {
		return err
	}

	output := bufio.NewWriter(os.Stdout)
	printer.Print(output, r)
	defer output.Flush()

	return nil
}

func CollectFileHashes(ctx context.Context, roots []string, algo hash.HashingAlgo) (map[string][]string, error) {
	pre := make(map[int64][]string, len(roots))

	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if d.Type()&os.ModeSymlink != 0 {
				return nil
			}

			fi, err := os.Stat(path)
			if err != nil {
				return err
			}

			size := fi.Size()
			if size == 0 {
				return nil
			}

			pre[size] = append(pre[size], path)

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	var p []string
	for _, slice := range pre {
		if len(slice) > 1 {
			p = append(p, slice...)
		}
	}

	results := make(map[string][]string, len(p))
	for _, root := range p {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			sum, err := hash.Hash(ctx, f, algo)
			if err != nil {
				return err
			}

			key := fmt.Sprintf("%x", sum)
			results[key] = append(results[key], path)

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func Run(ctx context.Context, cmd *cli.Command) (map[string][]string, error) {
	folders := []string{}
	if cmd.NArg() == 0 {
		folders = append(folders, ".")
	}

	for i := range cmd.NArg() {
		path := cmd.Args().Get(i)

		folders = append(folders, path)
	}

	algo, err := hash.ParseHashingAlgo(cmd.String("hash"))
	if err != nil {
		return nil, err
	}

	result, err := CollectFileHashes(ctx, folders, algo)
	if err != nil {
		return nil, err
	}

	return result, nil
}

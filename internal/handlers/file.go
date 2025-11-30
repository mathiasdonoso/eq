package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/mathiasdonoso/eq/pkg/hash"
	"github.com/urfave/cli/v3"
)

func CompareFiles(ctx context.Context, cmd *cli.Command) (map[string][]string, error) {
	if cmd.NArg() == 0 {
		return nil, fmt.Errorf("missing files to compare")
	}

	algo, err := hash.ParseHashingAlgo(cmd.String("hash"))
	if err != nil {
		return nil, err
	}

	files := make(map[string][]string, cmd.NArg())
	for i := range cmd.NArg() {
		filepath := cmd.Args().Get(i)

		fs, err := os.Lstat(filepath)
		if err != nil {
			return nil, err
		}

		if fs.IsDir() {
			continue
		}

		if fs.Mode() == os.ModeSymlink {
			continue
		}

		file, err := os.Open(filepath)
		if err != nil {
			return nil, err
		}

		b, err := hash.Hash(ctx, file, algo)
		if err != nil {
			return nil, err
		}

		sh := string(b)
		files[sh] = append(files[sh], filepath)
	}

	return files, nil
}

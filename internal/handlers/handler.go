package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"

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

			info, err := d.Info()
			if err != nil {
				return err
			}

			size := info.Size()
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

	var candidates []string
	for _, group := range pre {
		if len(group) > 1 {
			candidates = append(candidates, group...)
		}
	}

	if len(candidates) == 0 {
		return map[string][]string{}, nil
	}

	workerCount := runtime.NumCPU() * 2

	type job struct {
		path string
	}

	type result struct {
		hash string
		path string
		err  error
	}

	jobs := make(chan job, workerCount*2)
	results := make(chan result, workerCount*2)

	var wg sync.WaitGroup

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := range jobs {
				f, err := os.Open(j.path)
				if err != nil {
					results <- result{err: err}
					continue
				}

				sum, err := hash.Hash(ctx, f, algo)
				f.Close()
				if err != nil {
					results <- result{err: err}
					continue
				}

				results <- result{
					hash: fmt.Sprintf("%x", sum),
					path: j.path,
				}
			}
		}()
	}

	go func() {
		for _, p := range candidates {
			jobs <- job{path: p}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	final := make(map[string][]string)

	for res := range results {
		if res.err != nil {
			return nil, res.err
		}
		final[res.hash] = append(final[res.hash], res.path)
	}

	return final, nil
}

func Run(ctx context.Context, cmd *cli.Command) (map[string][]string, error) {
	folders := []string{}
	if cmd.NArg() == 0 {
		folders = append(folders, ".")
	} else {
		for i := range cmd.NArg() {
			folders = append(folders, cmd.Args().Get(i))
		}
	}

	algo, err := hash.ParseHashingAlgo(cmd.String("hash"))
	if err != nil {
		return nil, err
	}

	return CollectFileHashes(ctx, folders, algo)
}

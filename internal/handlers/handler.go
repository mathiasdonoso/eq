package handlers

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/mathiasdonoso/eq/pkg/hash"
	"github.com/urfave/cli/v3"
)

func Handler(ctx context.Context, cmd *cli.Command) error {
	r, err := Run(ctx, cmd)
	if err != nil {
		return err
	}

	r.Print()

	return nil
}

type Summary struct {
	Files map[string][]FileInfo
}

func (s Summary) Print() {
	for i := range s.Files {
		if len(s.Files[i]) == 1 {
			continue
		}

		fmt.Printf("Identical files:\n")
		for j := range s.Files[i] {
			fmt.Printf("Path: %s. Size: %d Bytes\n", s.Files[i][j].Path, s.Files[i][j].Size)
		}
		fmt.Println()
	}
}

type FileInfo struct {
	Path string
	Size int64
}

func CollectFileHashes(ctx context.Context, roots []string, algo hash.HashingAlgo) (Summary, error) {
	pre := make(map[int64][]FileInfo, len(roots))

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

			pre[size] = append(pre[size], FileInfo{
				Path: path,
				Size: size,
			})
			return nil
		})

		if err != nil {
			return Summary{}, err
		}
	}

	var candidates []FileInfo
	for _, group := range pre {
		if len(group) > 1 {
			candidates = append(candidates, group...)
		}
	}

	if len(candidates) == 0 {
		return Summary{}, nil
	}

	workerCount := runtime.NumCPU() * 2

	type job struct {
		path string
		size int64
	}

	type result struct {
		hash string
		path string
		size int64
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
					size: j.size,
				}
			}
		}()
	}

	go func() {
		for _, p := range candidates {
			jobs <- job{path: p.Path, size: p.Size}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	final := make(map[string][]FileInfo)

	for res := range results {
		if res.err != nil {
			return Summary{}, res.err
		}

		final[res.hash] = append(final[res.hash], FileInfo{
			Path: res.path,
			Size: res.size,
		})
	}

	s := Summary{
		Files: final,
	}

	return s, nil
}

func Run(ctx context.Context, cmd *cli.Command) (Summary, error) {
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
		return Summary{}, err
	}

	return CollectFileHashes(ctx, folders, algo)
}

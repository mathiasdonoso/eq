package hash

import (
	"context"
	"crypto/sha256"
	"io"

	"github.com/cespare/xxhash/v2"
	"github.com/zeebo/blake3"
)

type HashingFunc interface {
	io.Writer
	Sum(p []byte) []byte
}

func Hash(ctx context.Context, reader io.Reader, algorithm HashingAlgo) ([]byte, error) {
	var h HashingFunc

	switch algorithm {
	case SHA256:
		h = sha256.New()
	case BLAKE3:
		h = blake3.New()
	case XXH64:
		h = xxhash.New()
	default:
		h = blake3.New()
	}

	if _, err := io.Copy(h, reader); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

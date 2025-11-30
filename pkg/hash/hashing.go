package hash

import (
	"context"
	"crypto/sha256"
	"io"

	"github.com/cespare/xxhash/v2"
	"github.com/zeebo/blake3"
)

func Hash(ctx context.Context, reader io.Reader, algorithm HashingAlgo) ([]byte, error) {
	switch algorithm {
	case SHA256:
		return hashSHA256(reader)
	case BLAKE3:
		return hashBLAKE3(reader)
	case XXH64:
		return hashXXH64(reader)
	default:
		return hashBLAKE3(reader)
	}
}

func hashSHA256(reader io.Reader) ([]byte, error) {
	h := sha256.New()

	if _, err := io.Copy(h, reader); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func hashBLAKE3(reader io.Reader) ([]byte, error) {
	h := blake3.New()

	_, err := io.Copy(h, reader)
	if err != nil {
		return nil, err
	}

	sum := h.Sum(nil)
	return sum, nil
}

func hashXXH64(reader io.Reader) ([]byte, error) {
	h := xxhash.New()

	_, err := io.Copy(h, reader)
	if err != nil {
		return nil, err
	}

	sum := h.Sum(nil)
	return sum, nil
}

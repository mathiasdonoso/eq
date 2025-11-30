package hash

import (
	"fmt"
)

type HashingAlgo int

const (
	SHA256 HashingAlgo = iota
	BLAKE3
	XXH64
)

var algoToString = map[HashingAlgo]string{
	SHA256: "sha256",
	BLAKE3: "blake3",
	XXH64:  "xxh64",
}

var stringToAlgo = map[string]HashingAlgo{
	"sha256": SHA256,
	"blake3": BLAKE3,
	"xxh64":  XXH64,
}

func (h HashingAlgo) String() string {
	if v, ok := algoToString[h]; ok {
		return v
	}
	return fmt.Sprintf("unknown(%d)", int(h))
}

func ParseHashingAlgo(s string) (HashingAlgo, error) {
	if v, ok := stringToAlgo[s]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("invalid hashing algorithm: %s", s)
}

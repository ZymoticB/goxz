package filters

import (
	"errors"
	"github.com/ZymoticB/goxz/xz"
	"math"
)

var errInvalidLZMADictSize = errors.New("LZMA2 header contains invalid dictionary size")

type LZMADictSize int8

func (s LZMADictSize) size() (uint32, error) {
	size := int8(s)
	switch {
	case size > 40:
		return 0, errInvalidLZMADictSize
	case size == 40:
		return (4 * xz.GigaByte) - 1, nil
	case size%2 == 0:
		return uint32(math.Pow(2, float64((size/2)+12))), nil
	case size%2 == 1:
		return uint32(3 * (math.Pow(2, float64(((size-1)/2)+11)))), nil
	}
	return 0, errInvalidLZMADictSize
}

type LZMA2Header struct {
	DictSize LZMADictSize
}

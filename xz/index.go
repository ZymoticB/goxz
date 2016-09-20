package xz

import (
	"io"
)

type IndexIndicator byte

const (
	indexIndicator IndexIndicator = 0x00
)

type Index struct {
	Indicator IndexIndicator

	NumberOfRecords MultiByteInteger
	Records         []IndexRecord
	Padding         IndexPadding
	CRC32           CRC32
}

type IndexPadding [3]byte

type IndexRecord struct {
	UnpaddedSize     MultiByteInteger
	UncompressedSize MultiByteInteger
}

func (i *Index) read(r io.Reader) error {
	return nil
}

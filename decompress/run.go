package decompress

import (
	"errors"
	"log"
	"math"

	"github.com/ZymoticB/goxz/output"
)

const (
	KiloByte = 1024
	MegaByte = 1024 * KiloByte
	GigaByte = 1024 * MegaByte
)

func RunDecompress(source, dest string, out output) error {
	out.Print("RunDecompress")
	return nil
}

var errInvalidLZMADictSize = errors.New("LZMA2 header contains invalid dictionary size")

type LZMADictSize int8

func (LZMADictSize *s) size() (int32, error) {
	switch s {
	case s > 40:
		return -1, errInvalidLZMADictSize
	case s == 40:
		return (4 * GigaByte) - 1, nil
	case s%2 == 0:
		return 2 * *((s / 2) + 12), nil
	case s%2 == 1:
		return 3 * (2 * *(((s - 1) / 2) + 11)), nill
	}
}

type LZMA2Header struct {
	DictSize LZMADictSize
}

type CRC [4]byte

type XZStream struct {
	Header XZStreamHeader
	Blocks []XZBlock
	Index  XZIndex
	Footer XZStreamFooter
}

type XZStreamHeader struct {
	Magic [6]byte // FD 37 7A 58 5A 00
	Null  byte
	Flags byte
	CRC32 CRC
}

type XZStreamFooter struct {
	CRC32        CRC
	BackwardSize [4]byte
	Null         byte
	Flags        byte
	Magic        [2]byte
}

type XZStreamPadding struct {
	NullBytes [4]byte
}

type XZMultiByteInteger uint64 //1-9 byte variable length on disk encoding

type XZBlockIndexSharedHeader struct {
	BlockSizeSlashIndexIndicator byte
}

type FilterFlags struct {
	ID         XZMultiByteInteger
	Size       XZMultiByteInteger
	Properties []byte
}

type XZBlockHeader struct {
	// On-Disk Format is variable length so this structure cannot be
	// read directly from disk.

	XZBlockIndexSharedHeader

	Flags            byte
	CompressedSize   XZMultiByteInteger // Optional Field; if Flags & 0x40
	UncompressedSize XZMultiByteInteger // Optional Field; if Flags & 0x80
	FilterFlags      [4]FilterFlags     // Up to 4 based on Flags & 0x03
	Padding          []byte
	CRC32            CRC
}

type XZBlock struct {
	Header         XZBlockHeader
	CompressedData []byte
	Padding        []byte
	Check          []byte // variable length, sie and type depends on Stream Flags
}

type XZIndexRecord struct {
	UnpaddedSize     XZMultiByteInteger
	UncompressedSize XZMultiByteInteger
}

type XZIndex struct {
	XZBlockIndexSharedHeader // Always 0x00 for index

	NumberOfRecords XZMultiByteInteger
	Records         []XZIndexRecord
	Padding         []byte
	CRC32           CRC
}

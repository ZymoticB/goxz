package xz

import (
	"io"
)

type FilterFlags struct {
	ID         MultiByteInteger
	Size       MultiByteInteger
	Properties []byte
}

type BlockEncodedSize [1]byte

type BlockHeader struct {
	// On-Disk Format is variable length so this structure cannot be
	// read directly from disk.

	EncodedSize BlockEncodedSize

	Flags            byte
	CompressedSize   MultiByteInteger // Optional Field; if Flags & 0x40
	UncompressedSize MultiByteInteger // Optional Field; if Flags & 0x80
	FilterFlags      [4]FilterFlags   // Up to 4 based on Flags & 0x03
	Padding          []byte
	CRC32            CRC32
}

type Block struct {
	Header         BlockHeader
	CompressedData []byte
	Padding        []byte
	Check          []byte // variable length, sie and type depends on Stream Flags
}

func (b *Block) read(r io.Reader) error {
	return nil
}

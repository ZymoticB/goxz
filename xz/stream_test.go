package xz

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadStreamHeader(t *testing.T) {
	magic := [...]byte{0xFD, '7', 'z', 'X', 'Z', 0x00}
	magicR := bytes.NewReader(magic[:])

	flags := [...]byte{0x00, 0x1}
	flagR := bytes.NewReader(flags[:])

	crc := [...]byte{0x12, 0x13, 0x05, 0x72}
	crcAsInt := CRC32(0x72051312)
	crcR := bytes.NewReader(crc[:])

	r := io.MultiReader(magicR, flagR, crcR)

	var header StreamHeader

	err := header.read(r)
	assert.Nil(t, err)
	assert.Equal(t, header.Magic, StreamHeaderMagic(magic), "Byte order should match after reading from reader")

	assert.Equal(t, header.Flags, StreamFlags{0x0, 0x1}, "Flags should be read correctly")
	typ, err := header.Flags.getCheckType()
	assert.Nil(t, err)
	assert.Equal(t, typ, checkCRC32, "Check type should be computed correctly")
	assert.Equal(t, header.Flags.getCheckSize(), 4, "Check size should be computed correctly")
	assert.Equal(t, header.CRC, crcAsInt, "CRC should be read as little endian correctly")
}

func TestReadStreamFooter(t *testing.T) {
	crc := [...]byte{0x12, 0x13, 0x05, 0x72}
	crcAsInt := CRC32(0x72051312)
	crcR := bytes.NewReader(crc[:])

	bsize := [...]byte{0x01, 0x00, 0x00, 0x00}
	bsizeR := bytes.NewReader(bsize[:])

	flags := [...]byte{0x00, 0x1}
	flagR := bytes.NewReader(flags[:])

	magic := [...]byte{'Y', 'Z'}
	magicR := bytes.NewReader(magic[:])

	r := io.MultiReader(crcR, bsizeR, flagR, magicR)

	var footer StreamFooter

	err := footer.read(r)
	assert.Nil(t, err)
	assert.Equal(t, footer.CRC, crcAsInt, "CRC should be read from byte stream correctly")
	assert.Equal(t, footer.BackwardSize.getRealSize(), 8, "Backward size should be computed correctly")
	assert.Equal(t, footer.Flags, StreamFlags{0x0, 0x1}, "Flags should be read correctly")
	typ, err := footer.Flags.getCheckType()
	assert.Nil(t, err)
	assert.Equal(t, typ, checkCRC32, "Check type should be computed correctly")
	assert.Equal(t, footer.Flags.getCheckSize(), 4, "Check size should be computed correctly")
	assert.Equal(t, footer.Magic, StreamFooterMagic(magic), "Magic should be read correctly")
}

func TestReadNoPadding(t *testing.T) {
	padding := []byte{}
	r := bytes.NewReader(padding)
	var s Stream
	err := s.readPadding(r)

	assert.Nil(t, err)
	assert.Equal(t, len(s.Padding), 0, "there should be no padding in the stream")
}

func TestReadSomePadding(t *testing.T) {
	padding := []byte{0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00}
	r := bytes.NewReader(padding)
	var s Stream
	err := s.readPadding(r)

	assert.Nil(t, err)
	assert.Equal(t, len(s.Padding), 3, "should have read 2 blocks of padding")
}

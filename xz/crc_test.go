package xz

/*
echo "this is a test" | ../test/a.out
Bytes:  15
CRC-32: 0x72051312
CRC-64: 0x643D26FB7156AB08

a.out is the reference implementation from
http://tukaani.org/xz/xz-file-format.txt
*/

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRC32(t *testing.T) {
	//var table32 CRC32Table
	//table32.calculateLookupTable()
	// useful to generate lookup tables in a pasteable way

	testString := []byte("this is a test\n")

	crc32 := Crc32(testString, len(testString), 0)

	assert.Equal(t, crc32, uint32(0x72051312))
}

func TestCRC64(t *testing.T) {
	//var table64 CRC64Table
	//table64.init()
	// useful to generate lookup tables in a pasteable way

	testString := []byte("this is a test\n")

	crc64 := Crc64(testString, len(testString), 0)

	assert.Equal(t, crc64, uint64(0x643D26FB7156AB08))
}

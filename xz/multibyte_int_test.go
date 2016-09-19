package xz

import (
	"bytes"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

func TestReadSingleByteInt(t *testing.T) {
	var multibyteInt MultiByteInteger
	singleByte := []byte{0x1}

	reader := bytes.NewReader(singleByte)

	err := multibyteInt.Read(reader)
	assert.Equal(t, uint64(multibyteInt), uint64(1), "multibyte int should be read correctly")
	assert.Nil(t, err)

}

func TestReadMultiByteInt(t *testing.T) {
	var multibyteInt MultiByteInteger
	multibyteStream := []byte{0x81, 0x82, 0x83, 0x00}
	// 1 + 2 * 128 + 3 * 16384

	reader := bytes.NewReader(multibyteStream)

	err := multibyteInt.Read(reader)
	assert.Equal(t, uint64(multibyteInt), uint64(49409), "multibyte int should be read correctly")
	assert.Nil(t, err)

}

func TestDecodeSingleByteInt(t *testing.T) {
	var multibyteInt MultiByteInteger
	singleByte := []byte{0x1}

	err := multibyteInt.Decode(singleByte)
	assert.Equal(t, uint64(multibyteInt), uint64(1), "single byte int should be read correctly")
	assert.Nil(t, err)
}

func TestDecodeMultiByteInt(t *testing.T) {
	var multibyteInt MultiByteInteger
	multibyteStream := []byte{0x81, 0x82, 0x83, 0x00}
	// 1 + 2 * 128 + 3 * 16384

	err := multibyteInt.Decode(multibyteStream)
	assert.Equal(t, uint64(multibyteInt), uint64(49409), "multibyte int should be read correctly")
	assert.Nil(t, err)
}

func TestIntTooLong(t *testing.T) {
	var multibyteInt MultiByteInteger

	tooManyBytes := []byte{0xF1, 0xF2, 0xF3, 0xF4,
		0xF5, 0xF6, 0xF7, 0xF8,
		0xF9, 0xA}

	err := multibyteInt.Decode(tooManyBytes)
	assert.Equal(t, err, errMultiByteTooLong, "input too long should be rejected")
}

func TestEncodeSingleByte(t *testing.T) {
	var multibyteInt = MultiByteInteger(15)

	rawBytes, err := multibyteInt.Encode()

	assert.Nil(t, err)
	assert.Equal(t, rawBytes, []byte{0xF}, "single byte int should be encoded correctly")
}

func TestEncodingMultiByte(t *testing.T) {
	var multibyteInt = MultiByteInteger(128)

	rawBytes, err := multibyteInt.Encode()

	assert.Nil(t, err)
	assert.Equal(t, rawBytes, []byte{0x80, 0x1}, "multibyte int should be encoded correctly")
}

func TestRoundTrip(t *testing.T) {
	f := func(x uint64) bool {
		var multibyteIntSource = MultiByteInteger(x)
		var multibyteIntDest MultiByteInteger

		rawBytes, err := multibyteIntSource.Encode()

		if x > MultiByteMax {
			assert.Equal(t, err, errNumberTooLarge, "Numbers larger than 2^63-1 cannot be encoded")
			return true
		} else {
			assert.Nil(t, err)
		}

		err = multibyteIntDest.Decode(rawBytes)
		assert.Nil(t, err)
		return multibyteIntDest == multibyteIntSource
	}

	err := quick.Check(f, nil)
	assert.Nil(t, err)
}

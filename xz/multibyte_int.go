package xz

import (
	"errors"
	"io"
)

const (
	MultiByteMax = 1<<63 - 1
	MultiByteMin = 0
)

var errMultiByteTooLong = errors.New("MultiByte Integer Too Long")
var errNothingRead = errors.New("Non-EOF empty read??")
var errNumberTooLarge = errors.New("Unable to encode numbers > 2 ^ 63")

type MultiByteInteger uint64 //1-9 byte variable length on disk encoding

func (result *MultiByteInteger) Read(r io.Reader) error {
	buf := make([]byte, 9)
	i := 0
	for {
		var oneByte [1]byte
		numRead, err := r.Read(oneByte[:])
		if numRead != 1 {
			return errNothingRead
		}
		if err != nil {
			return err
		}
		buf[i] = oneByte[0]

		if oneByte[0]&0x80 == 0 {
			break
		}
		i += 1
	}
	return result.Decode(buf)
}

func (result *MultiByteInteger) Decode(buf []byte) error {
	var tmp uint64
	if len(buf) == 0 {
		*result = MultiByteInteger(0)
		return nil
	}

	if len(buf) > 9 {
		return errMultiByteTooLong
	}

	tmp = uint64(buf[0] & 0x7F)
	buf = buf[1:]

	for index, element := range buf {
		if element == 0x00 {
			break
		}

		tmp |= uint64(element&0x7F) << (uint32((index + 1) * 7))
	}
	*result = MultiByteInteger(tmp)
	return nil
}

func (source *MultiByteInteger) Encode() ([]byte, error) {
	num := uint64(*source)
	buf := make([]byte, 9)
	if num > MultiByteMax {
		return []byte{}, errNumberTooLarge
	}

	i := 0
	for num >= 0x80 {
		buf[i] = byte(num) | 0x80
		num = num >> 7
		i += 1
	}

	buf[i] = byte(num)
	return buf[:(i + 1)], nil
}

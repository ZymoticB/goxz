/*
This package notably does not support multi-stream files
*/
package xz

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

var errBadHeaderMagic = errors.New("Stream header has bad magic number")
var errBadFooterMagic = errors.New("Stream footer has bad magic number")
var errBadStreamFlags = errors.New("Stream flags first byte is not 0x00")
var errReservedFlagsUsed = errors.New("Reserved Stream Flags in use")

type Stream struct {
	Header  StreamHeader
	Blocks  []*Block
	Index   Index
	Footer  StreamFooter
	Padding []*StreamPadding
}

type StreamPadding [4]byte

type StreamFlags [2]byte
type StreamHeaderMagic [6]byte
type StreamFooterMagic [2]byte

//Constants, golang doesn't support non-basic type constants
var streamHeaderMagic = StreamHeaderMagic{0xFD, '7', 'z', 'X', 'Z', 0x00}
var streamFooterMagic = StreamFooterMagic{'Y', 'Z'}

type checkType string

const (
	checkNone   checkType = "none"
	checkCRC32  checkType = "crc32"
	checkCRC64  checkType = "crc64"
	checkSHA256 checkType = "sha256"
)

type streamContainerType string

const (
	isNone  streamContainerType = "none"
	isIndex streamContainerType = "index"
	isBlock streamContainerType = "block"
)

type StreamHeader struct {
	Magic StreamHeaderMagic // 0xFD 7 z X Z 0x00
	Flags StreamFlags       // 0x00 <flags>
	CRC   CRC32
}

type BackwardSize uint32
type StreamFooter struct {
	CRC          CRC32
	BackwardSize BackwardSize
	Flags        StreamFlags       // 0x00 <flags>
	Magic        StreamFooterMagic // YZ
}

func (stream *Stream) ReadStreamFromFile(path string) error {
	//open file
	//read stream
	return nil
}

func (stream *Stream) ReadStream(br *bufio.Reader) error {
	var err error
	err = stream.Header.read(br)
	if err != nil {
		return err
	}
	for {
		containerType, err := stream.readBlockOrIndex(br)
		if err != nil {
			return err
		}
		if containerType == isIndex {
			break
		}
	}
	err = stream.Footer.read(br)
	if err != nil {
		return err
	}
	err = stream.readPadding(br)
	if err != nil {
		return err
	}
	err = stream.validate()
	if err != nil {
		return err
	}
	return nil
}

func (s *Stream) readBlockOrIndex(br *bufio.Reader) (streamContainerType, error) {
	indicatorOrSize, err := br.Peek(1)
	if IndexIndicator(indicatorOrSize[0]) == indexIndicator {
		return isIndex, s.Index.read(br)
	} else {
		b := new(Block)
		err = b.read(br)
		if err != nil {
			return isNone, err
		}
		s.Blocks = append(s.Blocks, b)
		return isBlock, nil
	}
}

func (s *Stream) validate() error {
	return nil
}

func (s *Stream) readPadding(r io.Reader) error {
	for {
		p := new(StreamPadding)
		err := p.read(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		s.Padding = append(s.Padding, p)
	}
	return nil
}

func (padding *StreamPadding) read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, padding)
}

func (header *StreamHeader) read(r io.Reader) error {
	err := binary.Read(r, binary.BigEndian, &header.Magic)
	if err != nil {
		return err
	}
	if header.Magic != streamHeaderMagic {
		return errBadHeaderMagic
	}

	err = binary.Read(r, binary.BigEndian, &header.Flags)
	if err != nil {
		return err
	}
	if header.Flags[0] != 0x00 {
		return errBadStreamFlags
	}
	if header.Flags[1]&0xF0 != 0x0 {
		return errReservedFlagsUsed
	}

	err = binary.Read(r, binary.LittleEndian, &header.CRC)
	if err != nil {
		return err
	}
	return nil
}

func (footer *StreamFooter) read(r io.Reader) error {
	err := binary.Read(r, binary.LittleEndian, &footer.CRC)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.LittleEndian, &footer.BackwardSize)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &footer.Flags)
	if err != nil {
		return err
	}
	if footer.Flags[0] != 0x00 {
		return errBadStreamFlags
	}
	if footer.Flags[1]&0xF0 != 0x0 {
		return errReservedFlagsUsed
	}

	err = binary.Read(r, binary.BigEndian, &footer.Magic)
	if err != nil {
		return err
	}
	if footer.Magic != streamFooterMagic {
		return errBadFooterMagic
	}
	return nil
}

func (b *BackwardSize) getRealSize() int {
	return int((*b + 1) * 4)
}

func (flags *StreamFlags) getCheckType() (checkType, error) {
	flag := flags[1] & 0xF //only 4 bits of the second byte matter

	if flag == 0x0 {
		return checkNone, nil
	}

	if flag == 0x1 {
		return checkCRC32, nil
	}

	if flag == 0x4 {
		return checkCRC64, nil
	}

	if flag == 0xA {
		return checkSHA256, nil
	}

	return checkNone, errReservedFlagsUsed
}

func (flags *StreamFlags) getCheckSize() int {
	flag := flags[1] & 0xF //only 4 bits of the second byte matter

	if flag == 0x0 {
		return 0
	}

	if flag < 0x4 {
		return 4
	}

	if flag < 0x7 {
		return 8
	}

	if flag < 0xA {
		return 16
	}

	if flag < 0xD {
		return 32
	}
	return 64
}

package go_inthex

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
)

type RecordCode uint8

const (
	RecordData = RecordCode(iota)
	RecordEOF
	RecordExtendedSegmentAddress
	RecordStartSegmentAddress
	RecordExtendedLinearAddress
	RecordStartLinearAddress
)

type Record struct {
	Code    RecordCode
	Address uint16
	Data    []byte
}

func (r Record) String() string {
	if len(r.Data) > math.MaxUint8 {
		panic("data too long")
	}
	buf := make([]byte, 1+1+2+len(r.Data)+1)

	buf[0] = uint8(len(r.Data))
	binary.BigEndian.PutUint16(buf[1:], r.Address)
	buf[3] = uint8(r.Code)
	copy(buf[4:], r.Data)
	buf[len(buf)-1] = Checksum(buf[:len(buf)-1])
	return ":" + strings.ToUpper(hex.EncodeToString(buf)) + "\r\n"
}

func RecordFromString(data string) (Record, error) {
	if data[0] != ':' {
		return Record{}, fmt.Errorf("expected :, got %s", data[:1])
	}
	if len(data) < 1+(1+2+1)*2 {
		return Record{}, io.ErrUnexpectedEOF
	}

	b, err := hex.DecodeString(strings.Trim(data[1:], "\r\n"))
	if err != nil {
		return Record{}, err
	}

	if len(b) < 1 {
		return Record{}, io.ErrUnexpectedEOF
	}

	byteCount := b[0]
	if len(b) != 1+1+2+int(byteCount)+1 {
		return Record{}, errors.New("data mismatch")
	}

	var r Record

	r.Address = binary.BigEndian.Uint16(b[1:])
	r.Code = RecordCode(b[3])
	r.Data = b[4 : 4+int(byteCount)]

	checkSum := b[4+int(byteCount)]
	calculatedCksum := Checksum(b[:len(b)-1])
	if calculatedCksum != checkSum {
		return Record{}, errors.New("wrong checksum")
	}

	return r, nil
}

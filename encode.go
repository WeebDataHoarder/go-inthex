package go_inthex

import (
	"encoding/binary"
	"strings"
)

func Encode(stream *Stream) (hex []byte, err error) {

	var entries []string

	emitRecord := func(r Record) {
		entries = append(entries, r.String())
	}

	emitAddress := func(addr uint32) error {
		if addr <= 0x000FFFF0 && (addr&(^RecordExtendedSegmentAddressMask)) == 0 {
			emitRecord(Record{
				Code:    RecordExtendedSegmentAddress,
				Address: uint16(addr / 16),
				Data:    nil,
			})

			return nil
		} else {
			emitRecord(Record{
				Code:    RecordExtendedLinearAddress,
				Address: uint16(addr >> 16),
				Data:    nil,
			})

			return nil
		}
	}

	const recordSize = 16
	const sectionSize = 1 << 16

	for _, r := range stream.Regions {
		if err = emitAddress(r.Address); err != nil {
			return nil, err
		}

		// Add extra entries as needed
		for _, e := range r.Extra {
			entries = append(entries, e)
		}

		endAddress := r.Address + uint32(len(r.Data))

		for addr := r.Address; addr < endAddress; addr += sectionSize {
			if addr != r.Address {
				if err = emitAddress(addr); err != nil {
					return nil, err
				}
			}
			for subAddr := addr; subAddr < (addr+sectionSize) && subAddr < endAddress; subAddr += recordSize {
				writeSize := min(recordSize, endAddress-subAddr+1)
				dataOffset := subAddr - r.Address
				emitRecord(Record{
					Code:    RecordData,
					Address: uint16(subAddr & 0xFFFF),
					Data:    r.Data[dataOffset : dataOffset+writeSize],
				})
			}
		}
	}

	if stream.StartLinearAddress != 0 {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], stream.StartLinearAddress)
		emitRecord(Record{
			Code:    RecordStartLinearAddress,
			Address: 0,
			Data:    buf[:],
		})
	}

	if stream.StartSegmentAddress != 0 {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], stream.StartSegmentAddress)
		emitRecord(Record{
			Code:    RecordStartSegmentAddress,
			Address: 0,
			Data:    buf[:],
		})
	}
	emitRecord(Record{
		Code:    RecordEOF,
		Address: 0,
		Data:    nil,
	})

	return []byte(strings.Join(entries, "\r\n")), nil
}

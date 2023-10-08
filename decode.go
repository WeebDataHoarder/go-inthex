package go_inthex

import (
	"encoding/binary"
	"errors"
	"strings"
)

func Decode(data []byte) (stream *Stream, err error) {

	stream = &Stream{}

	var currentRegion Region
	var baseAddress uint32

	emitRegion := func() {
		if len(currentRegion.Data) > 0 || len(currentRegion.Extra) > 0 {
			stream.Regions = append(stream.Regions, currentRegion)
		}
		currentRegion = Region{}
	}

	for _, d := range strings.Split(string(data), "\n") {
		if len(d) == 0 {
			continue
		}

		if d[0] != ':' {
			currentRegion.Extra = append(currentRegion.Extra, strings.Trim(d, "\r\n"))
			continue
		}

		r, err := RecordFromString(d)
		if err != nil {
			return nil, err
		}

		switch r.Code {
		case RecordData:
			address := baseAddress + uint32(r.Address)
			if !currentRegion.IsContiguousAddress(address) {
				emitRegion()
				baseAddress = baseAddress + uint32(r.Address)
				currentRegion.Address = baseAddress
			}
			currentRegion.Append(r.Data)
		case RecordEOF:
			emitRegion()
			return stream, nil
		case RecordExtendedSegmentAddress:
			if len(r.Data) != 2 {
				return nil, errors.New("invalid ExtendedSegmentAddress length")
			}
			baseAddress = uint32(binary.BigEndian.Uint16(r.Data)) * 16
			if len(currentRegion.Data) == 0 {
				currentRegion.Address = baseAddress
			}
		case RecordStartSegmentAddress:
			if len(r.Data) != 4 {
				return nil, errors.New("invalid StartSegmentAddress length")
			}
			stream.StartLinearAddress = binary.BigEndian.Uint32(r.Data)
		case RecordExtendedLinearAddress:
			if len(r.Data) != 2 {
				return nil, errors.New("invalid ExtendedLinearAddress length")
			}
			baseAddress = uint32(binary.BigEndian.Uint16(r.Data)) << 16
			if len(currentRegion.Data) == 0 {
				currentRegion.Address = baseAddress
			}
		case RecordStartLinearAddress:
			if len(r.Data) != 4 {
				return nil, errors.New("invalid RecordStartLinearAddress length")
			}
			stream.StartLinearAddress = binary.BigEndian.Uint32(r.Data)
		}
	}

	return
}

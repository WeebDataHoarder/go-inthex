package go_inthex

type Region struct {
	Address uint32
	Data    []byte
	Extra   []string
}

func (r *Region) IsContiguousAddress(addr uint32) bool {
	return addr == (r.Address + uint32(len(r.Data)))
}

func (r *Region) Append(data []byte) {
	r.Data = append(r.Data, data...)
}

type Stream struct {
	Regions []Region

	StartSegmentAddress uint32
	StartLinearAddress  uint32
}

// Data Returns a binary image, from the lowest base address
func (s *Stream) Data() (data []byte, baseAddress uint32) {
	if len(s.Regions) == 0 {
		return nil, 0
	}
	baseAddress = s.Regions[0].Address
	topAddress := s.Regions[0].Address + uint32(len(s.Regions[0].Data))

	for _, r := range s.Regions[1:] {
		baseAddress = min(baseAddress, r.Address)
		topAddress = min(topAddress, r.Address+uint32(len(r.Data)))
	}

	data = make([]byte, topAddress-baseAddress)

	for _, r := range s.Regions {
		copy(data[r.Address:], r.Data)
	}

	return data, baseAddress
}

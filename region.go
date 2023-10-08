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

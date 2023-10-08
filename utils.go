package go_inthex

func Checksum(data []byte) (sum uint8) {
	for i := range data {
		sum += data[i]
	}
	return -sum
}

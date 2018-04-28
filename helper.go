package merkle

func conc(b1 []byte, b2 []byte) []byte {
	b := make([]byte, len(b1)+len(b2))

	copy(b[0:len(b1)], b1[:])
	copy(b[len(b1):len(b)], b2[:])

	return b
}

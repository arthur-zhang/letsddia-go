package bitset

type BitSet []byte

func New(size uint32) BitSet {
	size = (size + 7) &^ 7
	return make([]byte, size>>3)
}
func (bs *BitSet) Len() uint32 {
	return uint32(len(*bs)) << 3
}

func (bs BitSet) Set(n uint32) {
	if n >= bs.Len() {
		panic("out of range")
	}
	bs[n>>3] |= 1 << (n & 7)
}
func (bs BitSet) IsSet(n uint32) bool {
	if n >= bs.Len() {
		panic("out of range")
	}
	return bs[n>>3]&(1<<(n&7)) != 0
}

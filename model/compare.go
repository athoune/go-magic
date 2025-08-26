package model

type Compare struct {
	Not         bool // !
	X           bool // special value, always return true
	Endianness  byte // n, e, b <- native, little, big
	Type        *Type
	Operation   byte // = > < & ^ ~
	StringValue string
	FloatValue  float64
	IntValue    int64
	UIntValue   uint64
	QuadValue   []int64
	UQuadValue  []uint64
}

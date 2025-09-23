package model

type Value struct {
	Family         TypeFamily
	StringValue    string
	FloatValue     float64
	IntValue       int64
	UIntValue      uint64
	RawBinaryValue uint64
}

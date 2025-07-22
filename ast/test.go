package ast

type Type int

const (
	Byte Type = iota
	Short
	Long
	Quad
)

type Offset struct {
	Level     int
	Value     int64
	Dynamic   bool
	DynOffset int64
	DynType   byte
	DynAction byte
	DynArg    int64
}

type Test struct {
	Offset   Offset
	Type     string
	SubTests []Test
}

type Compare struct {
	Operation   byte // = > < & ^ ~
	Not         bool // !
	StringValue string
	FloatValue  float64
	IntValue    int64
	Type        byte // s f i
}

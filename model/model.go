package model

type Offset struct {
	Level       int
	Relative    bool
	Value       int64
	Dynamic     bool
	DynOffset   int64
	DynType     byte
	DynOperator byte
	DynArg      int64
}

type Compare struct {
	Operation   byte // = > < & ^ ~
	Not         bool // !
	StringValue string
	FloatValue  float64
	IntValue    int64
	QuadValue   []int64
	Type        Clue
}

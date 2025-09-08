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

package model

type Compare struct {
	Not         bool // !
	X           bool // special value, always return true
	Type        *Type
	Comparator  byte // = > < & ^ ~
	RawExpected string
	Expected    *Value
}

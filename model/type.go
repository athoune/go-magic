package model

const (
	TYPE_CLUE_INT = Clue(iota)
	TYPE_CLUE_UINT
	TYPE_CLUE_FLOAT
	TYPE_CLUE_STRING
)

type Clue int
type StringOptions uint16

const (
	STRING_OPTIONS_NONE                   = StringOptions(0b0)
	STRING_OPTIONS_COMPACT_WITH_SPACES    = StringOptions(0b1 << 1)
	STRING_OPTIONS_FULL_WORD              = StringOptions(0b1 << 2)
	STRING_OPTIONS_CASE_INSENSITIVE_UPPER = StringOptions(0b1 << 3)
	STRING_OPTIONS_CASE_INSENSITIVE_LOWER = StringOptions(0b1 << 4)
	STRING_OPTIONS_TEXT_FILE              = StringOptions(0b1 << 5)
	STRING_OPTIONS_BINARY_FILE            = StringOptions(0b1 << 6)
	STRING_OPTIONS_TRIMMED                = StringOptions(0b1 << 7)
	REGEX_OPTIONS_OFFSET_START            = StringOptions(0b1 << 8)
)

type Type struct {
	Root      string // ubelong -> long
	ByteOrder BYTE_ORDER
	Name      string
	Clue      Clue
	// Filter
	FilterOperator       byte
	FilterBinaryArgument uint64
	FilterStringArgument string
	// when root == "string" or "search"
	StringOptions StringOptions
	// when root == "search"
	SearchRange int
	SearchCount int
}

var Types map[string]Clue

func init() {
	// [FIXME] use type shortener
	Types = map[string]Clue{
		"byte":        TYPE_CLUE_INT,
		"ubyte":       TYPE_CLUE_INT,
		"short":       TYPE_CLUE_INT,
		"leshort":     TYPE_CLUE_INT,
		"beshort":     TYPE_CLUE_INT,
		"ushort":      TYPE_CLUE_INT,
		"ubeshort":    TYPE_CLUE_INT,
		"uleshort":    TYPE_CLUE_INT,
		"u4":          TYPE_CLUE_INT,
		"long":        TYPE_CLUE_INT,
		"belong":      TYPE_CLUE_INT,
		"ulong":       TYPE_CLUE_INT,
		"ubelong":     TYPE_CLUE_INT,
		"ulelong":     TYPE_CLUE_INT,
		"lelong":      TYPE_CLUE_INT,
		"u8":          TYPE_CLUE_INT,
		"quad":        TYPE_CLUE_INT,
		"ubequad":     TYPE_CLUE_INT,
		"uquad":       TYPE_CLUE_INT,
		"ulequad":     TYPE_CLUE_INT,
		"float":       TYPE_CLUE_FLOAT,
		"double":      TYPE_CLUE_FLOAT,
		"string":      TYPE_CLUE_STRING,
		"pstring":     TYPE_CLUE_STRING,
		"ustring":     TYPE_CLUE_STRING,
		"date":        TYPE_CLUE_STRING,
		"lemsdosdate": TYPE_CLUE_STRING,
		"qdate":       TYPE_CLUE_STRING,
		"ldate":       TYPE_CLUE_STRING,
		"qldate":      TYPE_CLUE_STRING,
		"qwdate":      TYPE_CLUE_STRING,
		"uledate":     TYPE_CLUE_STRING,
		"ubeqdate":    TYPE_CLUE_STRING,
		"lemsdostime": TYPE_CLUE_STRING,
		"beid3":       TYPE_CLUE_STRING,
		"bequad":      TYPE_CLUE_INT,
		"befloat":     TYPE_CLUE_FLOAT,
		"bedouble":    TYPE_CLUE_FLOAT,
		"bedate":      TYPE_CLUE_STRING,
		"beqdate":     TYPE_CLUE_STRING,
		"beldate":     TYPE_CLUE_STRING,
		"beqldate":    TYPE_CLUE_STRING,
		"beqwdate":    TYPE_CLUE_STRING,
		"bestring16":  TYPE_CLUE_STRING,
		"leid3":       TYPE_CLUE_STRING,
		"lequad":      TYPE_CLUE_INT,
		"lefloat":     TYPE_CLUE_FLOAT,
		"ledouble":    TYPE_CLUE_FLOAT,
		"ledate":      TYPE_CLUE_STRING,
		"leqdate":     TYPE_CLUE_STRING,
		"leldate":     TYPE_CLUE_STRING,
		"leqldate":    TYPE_CLUE_STRING,
		"leqwdate":    TYPE_CLUE_STRING,
		"lestring16":  TYPE_CLUE_STRING,
		"melong":      TYPE_CLUE_STRING,
		"medate":      TYPE_CLUE_STRING,
		"meldate":     TYPE_CLUE_STRING,
		"indirect":    TYPE_CLUE_STRING,
		"name":        TYPE_CLUE_STRING,
		"use":         TYPE_CLUE_STRING,
		"regex":       TYPE_CLUE_STRING,
		"search":      TYPE_CLUE_STRING,
		"default":     TYPE_CLUE_STRING,
		"clear":       TYPE_CLUE_STRING,
		"der":         TYPE_CLUE_STRING,
		"guid":        TYPE_CLUE_STRING,
		"offset":      TYPE_CLUE_STRING,
	}
}

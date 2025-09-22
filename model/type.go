package model

const (
	TYPE_CLUE_INT = Family(iota)
	TYPE_CLUE_UINT
	TYPE_CLUE_FLOAT
	TYPE_CLUE_STRING
)

type Family int
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
	Family    Family
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

var Types map[string]Family

func init() {
	// [FIXME] use type shortener
	Types = map[string]Family{
		"byte":        TYPE_CLUE_INT,
		"ubyte":       TYPE_CLUE_UINT,
		"short":       TYPE_CLUE_INT,
		"ushort":      TYPE_CLUE_UINT,
		"u4":          TYPE_CLUE_INT,
		"long":        TYPE_CLUE_INT,
		"ulong":       TYPE_CLUE_UINT,
		"u8":          TYPE_CLUE_INT,
		"quad":        TYPE_CLUE_INT,
		"uquad":       TYPE_CLUE_UINT,
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
		"id3":         TYPE_CLUE_STRING,
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

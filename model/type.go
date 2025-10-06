package model

const (
	TYPE_FAMILY_INT = TypeFamily(iota)
	TYPE_FAMILY_UINT
	TYPE_FAMILY_FLOAT
	TYPE_FAMILY_STRING
)

type TypeFamily int
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
	Root       string // ubelong -> long
	ByteOrder  BYTE_ORDER
	Name       string
	TypeFamily TypeFamily
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

var Types map[string]TypeFamily

func init() {
	// [FIXME] use type shortener
	Types = map[string]TypeFamily{
		"byte":        TYPE_FAMILY_INT,
		"short":       TYPE_FAMILY_INT,
		"u4":          TYPE_FAMILY_INT,
		"4":           TYPE_FAMILY_INT,
		"long":        TYPE_FAMILY_INT,
		"u8":          TYPE_FAMILY_INT,
		"8":           TYPE_FAMILY_INT,
		"quad":        TYPE_FAMILY_INT,
		"float":       TYPE_FAMILY_FLOAT,
		"double":      TYPE_FAMILY_FLOAT,
		"string":      TYPE_FAMILY_STRING,
		"string16":    TYPE_FAMILY_STRING,
		"pstring":     TYPE_FAMILY_STRING,
		"ustring":     TYPE_FAMILY_STRING,
		"date":        TYPE_FAMILY_STRING,
		"msdosdate":   TYPE_FAMILY_STRING,
		"msdostime":   TYPE_FAMILY_STRING,
		"qdate":       TYPE_FAMILY_STRING,
		"ldate":       TYPE_FAMILY_STRING,
		"qldate":      TYPE_FAMILY_STRING,
		"qwdate":      TYPE_FAMILY_STRING,
		"uledate":     TYPE_FAMILY_STRING,
		"ubeqdate":    TYPE_FAMILY_STRING,
		"lemsdostime": TYPE_FAMILY_STRING,
		"beid3":       TYPE_FAMILY_STRING,
		"id3":         TYPE_FAMILY_STRING,
		"meldate":     TYPE_FAMILY_STRING,
		"indirect":    TYPE_FAMILY_STRING,
		"name":        TYPE_FAMILY_STRING,
		"use":         TYPE_FAMILY_STRING,
		"regex":       TYPE_FAMILY_STRING,
		"search":      TYPE_FAMILY_STRING,
		"default":     TYPE_FAMILY_STRING,
		"clear":       TYPE_FAMILY_STRING,
		"der":         TYPE_FAMILY_STRING,
		"guid":        TYPE_FAMILY_STRING,
		"offset":      TYPE_FAMILY_STRING,
	}
}

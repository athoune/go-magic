package ast

const (
	TYPE_CLUE_INT = Clue(iota)
	TYPE_CLUE_FLOAT
	TYPE_CLUE_STRING
)

type Clue int
type Type struct {
	Name  string
	Clue_ Clue
}

var Types map[string]*Type

func init() {
	Types = map[string]*Type{
		"byte": {
			Clue_: TYPE_CLUE_INT,
		},
		"ubyte": {
			Clue_: TYPE_CLUE_INT,
		},
		"short": {
			Clue_: TYPE_CLUE_INT,
		},
		"ubeshort": {
			Clue_: TYPE_CLUE_INT,
		},
		"long": {
			Clue_: TYPE_CLUE_INT,
		},
		"quad": {
			Clue_: TYPE_CLUE_INT,
		},
		"float": {
			Clue_: TYPE_CLUE_FLOAT,
		},
		"double": {
			Clue_: TYPE_CLUE_FLOAT,
		},
		"string": {
			Clue_: TYPE_CLUE_STRING,
		},
		"pstring": {
			Clue_: TYPE_CLUE_STRING,
		},
		"date": {
			Clue_: TYPE_CLUE_STRING,
		},
		"qdate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"ldate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"qldate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"qwdate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"beid3": {
			Clue_: TYPE_CLUE_STRING,
		},
		"beshort": {
			Clue_: TYPE_CLUE_INT,
		},
		"belong": {
			Clue_: TYPE_CLUE_INT,
		},
		"bequad": {
			Clue_: TYPE_CLUE_INT,
		},
		"befloat": {
			Clue_: TYPE_CLUE_FLOAT,
		},
		"bedouble": {
			Clue_: TYPE_CLUE_FLOAT,
		},
		"bedate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"beqdate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"beldate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"beqldate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"beqwdate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"bestring16": {
			Clue_: TYPE_CLUE_STRING,
		},
		"leid3": {
			Clue_: TYPE_CLUE_STRING,
		},
		"leshort": {
			Clue_: TYPE_CLUE_INT,
		},
		"lelong": {
			Clue_: TYPE_CLUE_INT,
		},
		"lequad": {
			Clue_: TYPE_CLUE_INT,
		},
		"lefloat": {
			Clue_: TYPE_CLUE_FLOAT,
		},
		"ledouble": {
			Clue_: TYPE_CLUE_FLOAT,
		},
		"ledate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"leqdate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"leldate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"leqldate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"leqwdate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"lestring16": {
			Clue_: TYPE_CLUE_STRING,
		},
		"melong": {
			Clue_: TYPE_CLUE_STRING,
		},
		"medate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"meldate": {
			Clue_: TYPE_CLUE_STRING,
		},
		"indirect": {
			Clue_: TYPE_CLUE_STRING,
		},
		"name": {
			Clue_: TYPE_CLUE_STRING,
		},
		"use": {
			Clue_: TYPE_CLUE_STRING,
		},
		"regex": {
			Clue_: TYPE_CLUE_STRING,
		},
		"search": {
			Clue_: TYPE_CLUE_STRING,
		},
		"default": {
			Clue_: TYPE_CLUE_STRING,
		},
		"clear": {
			Clue_: TYPE_CLUE_STRING,
		},
		"der": {
			Clue_: TYPE_CLUE_STRING,
		},
		"guid": {
			Clue_: TYPE_CLUE_STRING,
		},
		"offset": {
			Clue_: TYPE_CLUE_STRING,
		},
	}
	for k, v := range Types {
		v.Name = k
	}
}

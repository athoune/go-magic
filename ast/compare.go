package ast

import (
	"encoding/binary"

	"github.com/athoune/go-magic/model"
)

var TYPE_COMPARATORS map[string]typeFilter

func init() {
	TYPE_COMPARATORS = map[string]typeFilter{
		"byte":   &ByteComparator{},
		"short":  &ShortComparator{},
		"long":   &LongComparator{},
		"quad":   &QuadComparator{},
		"string": &StringComparator{},
	}
}

type typeFilter interface {
	compare(value []byte, bo binary.ByteOrder, signed model.SIGN,
		typ *model.Type, compare *model.Compare) (bool, bool, error)
}

type ByteComparator struct {
}

func (l *ByteComparator) compare(buff []byte, bo binary.ByteOrder, signed model.SIGN,
	typ *model.Type, compare *model.Compare) (bool, bool, error) {
	if signed == model.SIGNED {
		return operation(
				compare.IntValue,
				int64(buff[0]),
				compare.Operation),
			false,
			nil
	} else {
		return operation(
				compare.UIntValue,
				uint64(buff[0]),
				compare.Operation),
			false,
			nil
	}
}

type ShortComparator struct {
}

func (l *ShortComparator) compare(buff []byte, bo binary.ByteOrder, signed model.SIGN,
	typ *model.Type, compare *model.Compare) (bool, bool, error) {
	if signed == model.SIGNED {
		v, err := filterInt(int64(bo.Uint16(buff)), typ)
		if err != nil {
			return false, false, err
		}
		return operation(
				compare.IntValue,
				v,
				compare.Operation),
			false,
			nil
	} else {
		v, err := filterUInt(uint64(bo.Uint16(buff)), typ)
		if err != nil {
			return false, false, err
		}
		return operation(
				compare.UIntValue,
				v,
				compare.Operation),
			false,
			nil
	}
}

type LongComparator struct {
}

func (l *LongComparator) compare(buff []byte, bo binary.ByteOrder, signed model.SIGN,
	typ *model.Type, compare *model.Compare) (bool, bool, error) {
	if signed == model.SIGNED {
		v, err := filterInt(int64(bo.Uint32(buff)), typ)
		if err != nil {
			return false, false, err
		}
		return operation(
				compare.IntValue,
				v,
				compare.Operation),
			false,
			nil
	} else {
		v, err := filterUInt(uint64(bo.Uint32(buff)), typ)
		if err != nil {
			return false, false, err
		}
		return operation(
				compare.UIntValue,
				v,
				compare.Operation),
			false,
			nil
	}
}

type QuadComparator struct {
}

func (l *QuadComparator) compare(buff []byte, bo binary.ByteOrder, signed model.SIGN,
	typ *model.Type, compare *model.Compare) (bool, bool, error) {
	if signed == model.SIGNED {
		v, err := filterInt(int64(bo.Uint32(buff)), typ)
		if err != nil {
			return false, false, err
		}
		return operation(
				compare.IntValue,
				v,
				compare.Operation),
			false,
			nil
	} else {
		v, err := filterUInt(uint64(bo.Uint32(buff)), typ)
		if err != nil {
			return false, false, err
		}
		return operation(
				compare.UIntValue,
				v,
				compare.Operation),
			false,
			nil
	}
}

type StringComparator struct {
}

func (l *StringComparator) compare(buff []byte, bo binary.ByteOrder, signed model.SIGN,
	typ *model.Type, compare *model.Compare) (bool, bool, error) {

	return operationString(compare.StringValue,
			string(buff),
			compare.Operation),
		false,
		nil

}

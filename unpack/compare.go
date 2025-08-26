package unpack

import (
	"encoding/binary"

	"github.com/athoune/go-magic/model"
)

type typeComparator func(value []byte, bo binary.ByteOrder, signed model.SIGN,
	typ *model.Type, compare *model.Compare) (bool, bool, error)

var TYPE_COMPARATORS map[string]typeComparator

func init() {
	TYPE_COMPARATORS = map[string]typeComparator{
		"byte":   ByteComparator,
		"short":  ShortComparator,
		"long":   LongComparator,
		"quad":   QuadComparator,
		"string": StringComparator,
	}
}

func ByteComparator(buff []byte, bo binary.ByteOrder, signed model.SIGN,
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

func ShortComparator(buff []byte, bo binary.ByteOrder, signed model.SIGN,
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

func LongComparator(buff []byte, bo binary.ByteOrder, signed model.SIGN,
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

func QuadComparator(buff []byte, bo binary.ByteOrder, signed model.SIGN,
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

func StringComparator(buff []byte, bo binary.ByteOrder, signed model.SIGN,
	typ *model.Type, compare *model.Compare) (bool, bool, error) {

	return operationString(compare.StringValue,
			string(buff),
			compare.Operation),
		false,
		nil

}

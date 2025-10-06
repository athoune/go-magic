package ast

import (
	"errors"
	"fmt"

	"github.com/athoune/go-magic/model"
)

func FilterValue(typ *model.Type, value *model.Value) error {
	if typ.FilterOperator == ' ' { // there is no filter
		return nil
	}
	switch typ.TypeFamily {
	case model.TYPE_FAMILY_STRING:
		return errors.New("string filters are not implemented")
	case model.TYPE_FAMILY_FLOAT:
		return errors.New("float filters are not implemented")
	case model.TYPE_FAMILY_INT, model.TYPE_FAMILY_UINT:
		vv := value.RawBinaryValue
		switch typ.FilterOperator {
		case '&':
			vv &= typ.FilterBinaryArgument
		case '|':
			vv |= typ.FilterBinaryArgument
		case '^':
			vv ^= typ.FilterBinaryArgument
		/*
			case '/':
			case '%':
			case '+':
			case '-':
			case '*':
		*/
		default:
			return fmt.Errorf("unknown operator : %v", typ.FilterOperator)
		}
		if typ.TypeFamily == model.TYPE_FAMILY_INT {
			value.IntValue = int64(vv)
		} else {
			value.UIntValue = vv
		}
	}

	return nil
}

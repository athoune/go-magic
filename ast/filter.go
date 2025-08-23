package ast

import (
	"fmt"
	"strconv"

	"github.com/athoune/go-magic/model"
)

func filterInt(arg int64, typ *model.Type) (int64, error) {
	if typ.Operator == 0 {
		return arg, nil
	}
	r, err := strconv.ParseInt(typ.Arg, 0, 64)
	if err != nil {
		return 0, err
	}
	switch typ.Operator {
	case '&':
		return arg & r, nil
	}
	return 0, fmt.Errorf("unknown operator: %v", typ.Operator)
}

func filterUInt(arg uint64, typ *model.Type) (uint64, error) {
	if typ.Operator == 0 {
		return arg, nil
	}
	r, err := strconv.ParseInt(typ.Arg, 0, 64)
	if err != nil {
		return 0, err
	}
	switch typ.Operator {
	case '&':
		return arg & uint64(r), nil
	}
	return 0, fmt.Errorf("unknown operator: %v", typ.Operator)
}

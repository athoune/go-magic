package ast

import (
	"fmt"
	"io"

	"github.com/athoune/go-magic/model"
)

func ReadAndCompareValues(actual *model.Type, expected *model.Value, data io.Reader) (bool, error) {
	if actual.TypeFamily != expected.Family {
		return false,
			fmt.Errorf("you can't compare different value type : %v vs %v",
				actual.TypeFamily, expected.Family)
	}
	/*
		actualValue, _, err := unpack.ReadValue(actual, data)
		if err != nil {
			return false, err
		}
	*/

	return false, nil
}

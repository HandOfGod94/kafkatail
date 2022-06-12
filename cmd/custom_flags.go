package cmd

import (
	"fmt"
	"strconv"

	"github.com/samber/mo"
)

type PartitionFlag struct {
	value mo.Either[int, string]
}

func (ef *PartitionFlag) Type() string {
	return "either[int, string]"
}

func (ef *PartitionFlag) Set(arg string) error {
	res, err := strconv.Atoi(arg)
	if err != nil {
		if arg != "all" {
			return fmt.Errorf(`expected int (0,1,2,...) or string "all"`)
		}

		ef.value = mo.Right[int](arg)
		return nil
	}

	ef.value = mo.Left[int, string](res)
	return nil
}

func (ef *PartitionFlag) String() string {
	if ef.value.IsLeft() {
		return fmt.Sprintf("%v", ef.value.MustLeft())
	}

	if ef.value.IsRight() {
		return ef.value.MustRight()
	}

	return "<empty>"
}

func (ef *PartitionFlag) Get() mo.Either[int, string] {
	return ef.value
}

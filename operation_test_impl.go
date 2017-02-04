package patchstructure

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/pointerstructure"
)

// RFC6902 4.6
func opTest(op *Operation, v interface{}) (interface{}, error) {
	// Parse the path
	pointer, err := pointerstructure.Parse(op.Path)
	if err != nil {
		return v, err
	}

	// Target location must exist
	target, err := pointer.Get(v)
	if err != nil {
		return v, err
	}

	// Perform the test with "reflect". This can be improved in the future
	// with good reason since there are subtle type issues with reflect
	// DeepEqual.
	err = nil
	if !reflect.DeepEqual(target, op.Value) {
		err = fmt.Errorf("values not equal: %#v != %#v", target, op.Value)
	}

	return v, err
}

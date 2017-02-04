package patchstructure

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/pointerstructure"
)

// RFC6902 4.4
func opMove(op *Operation, v interface{}) (interface{}, error) {
	// Parse the path. We do this even though we don't use it to
	// avoid syntax errors causing partial applies.
	to, err := pointerstructure.Parse(op.Path)
	if err != nil {
		return v, err
	}

	// Parse the from path, which must exist
	from, err := pointerstructure.Parse(op.From)
	if err != nil {
		return v, err
	}

	// "The "from" location MUST NOT be a proper prefix of the "path"
	// location; i.e., a location cannot be moved into one of its children."
	if len(from.Parts) < len(to.Parts) {
		if reflect.DeepEqual(from.Parts, to.Parts[:len(from.Parts)]) {
			return v, fmt.Errorf(
				"move cannot move into a child path of the from path")
		}
	}

	// Get the from value, which must exist
	fromValue, err := from.Get(v)
	if err != nil {
		return v, err
	}

	// "This operation is functionally identical to a "remove" operation on
	// the "from" location, followed immediately by an "add" operation at
	// the target location with the value that was just removed."
	removeOp := &Operation{
		Op:   OpRemove,
		Path: op.From,
	}

	addOp := &Operation{
		Op:    OpAdd,
		Path:  op.Path,
		Value: fromValue,
	}

	// Remove first
	v, err = removeOp.Apply(v)
	if err != nil {
		return v, err
	}

	// Add
	return addOp.Apply(v)
}

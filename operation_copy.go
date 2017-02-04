package patchstructure

import (
	"github.com/mitchellh/pointerstructure"
)

// RFC6902 4.5
func opCopy(op *Operation, v interface{}) (interface{}, error) {
	// Parse the path. We do this even though we don't use it to
	// avoid syntax errors causing partial applies.
	_, err := pointerstructure.Parse(op.Path)
	if err != nil {
		return v, err
	}

	// Parse the from path, which must exist
	from, err := pointerstructure.Parse(op.From)
	if err != nil {
		return v, err
	}

	// Get the from value, which must exist
	fromValue, err := from.Get(v)
	if err != nil {
		return v, err
	}

	// "This operation is functionally identical to an "add" operation at the
	// target location using the value specified in the "from" member."

	addOp := &Operation{
		Op:    OpAdd,
		Path:  op.Path,
		Value: fromValue,
	}

	// Add
	return addOp.Apply(v)
}

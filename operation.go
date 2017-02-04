package patchstructure

import (
	"fmt"
)

// Operation represents a single operation to apply to a structure.
//
// Note that Value and From are dependent on the operation. Please see
// the JSON patch documentation for details on this since the semantics for
// the Go patch is identical.
type Operation struct {
	Op    Op          // Op is the operation type to apply
	Path  string      // Path is required
	Value interface{} // Optional depending on op
	From  string      // Optional depending on op
}

// Op is an enum representing the supported operations for a patch.
//
// The values should obviously match the JSON patch operations and their
// semantics are meant to be identical for Go structures.
type Op int

const (
	OpInvalid Op = iota // Set zero to invalid to prevent accidental adds
	OpAdd
	OpRemove
	OpReplace
	OpMove
	OpCopy
	OpTest
)

// String format of an operation matching what it should be if JSON encoded.
func (o Op) String() string {
	return opString[o]
}

// Apply performs the operation on the value v. The value v will be modified.
// In the case of an error, v may still be modified. If you wish to protect
// against partial failure, please deep copy the object prior to changes.
func (o *Operation) Apply(v interface{}) (interface{}, error) {
	f, ok := opApplyMap[o.Op]
	if !ok {
		return nil, fmt.Errorf("unknown operation: %s", o.Op)
	}

	result, err := f(o, v)
	if err != nil {
		return nil, fmt.Errorf("error applying operation %s: %s", o.Op, err)
	}

	return result, nil
}

var opString = map[Op]string{
	OpInvalid: "invalid",
	OpAdd:     "add",
	OpRemove:  "remove",
	OpReplace: "replace",
	OpMove:    "move",
	OpCopy:    "copy",
	OpTest:    "test",
}

// onApplyFunc is the type used internally for applying operations.
type opApplyFunc func(*Operation, interface{}) (interface{}, error)

// onApplyMap is the map used for lookup for the action to perform
// when applying an operation.
var opApplyMap = map[Op]opApplyFunc{
	OpAdd: opAdd,
}

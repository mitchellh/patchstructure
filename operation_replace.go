package patchstructure

import (
	"github.com/mitchellh/pointerstructure"
)

// RFC6902 4.3
func opReplace(op *Operation, v interface{}) (interface{}, error) {
	// Parse the path
	pointer, err := pointerstructure.Parse(op.Path)
	if err != nil {
		return v, err
	}

	// The only thing we need to check is that the pointer path actually
	// exists. If it doesn't, it is an error. To quote the RFC:
	//
	// "The target location MUST exist for the operation to be successful."
	if _, err := pointer.Get(v); err != nil {
		return v, err
	}

	// Set always does the right thing
	return pointer.Set(v, op.Value)
}

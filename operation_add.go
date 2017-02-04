package patchstructure

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/pointerstructure"
)

// RFC6902 4.1
func opAdd(op *Operation, v interface{}) (interface{}, error) {
	// Parse the path
	pointer, err := pointerstructure.Parse(op.Path)
	if err != nil {
		return nil, err
	}

	// If the pointer is root, then we apply directly to it since it'll
	// replace the entire doc. RFC quote below.
	if pointer.IsRoot() {
		// "The root of the target document - whereupon the specified value
		//  becomes the entire content of the target document."
		return pointer.Set(v, op.Value)
	}

	// Get the path that we want to add to (the parent)
	parent, err := pointer.Parent().Get(v)
	if err != nil {
		return nil, err
	}

	// The type will determine how we handle this
	parentVal := reflect.ValueOf(parent)
	switch parentVal.Kind() {
	case reflect.Map:
		// "If the target location specifies an object member that does not
		// already exist, a new member is added to the object."
		//
		// "If the target location specifies an object member that does exist,
		// that member's value is replaced."
		return pointer.Set(v, op.Value)

	default:
		return nil, fmt.Errorf(
			"can only add to maps, slices, arrays, or structs, got %q",
			parentVal.Kind())
	}
}

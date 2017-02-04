package patchstructure

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/mitchellh/pointerstructure"
)

// RFC6902 4.1
func opAdd(op *Operation, v interface{}) (interface{}, error) {
	// Parse the path
	pointer, err := pointerstructure.Parse(op.Path)
	if err != nil {
		return v, err
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
		return v, err
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

	case reflect.Slice:
		return opAddSlice(pointer, parentVal, op, v)

	default:
		return v, fmt.Errorf(
			"can only add to maps, slices, arrays, or structs, got %q",
			parentVal.Kind())
	}
}

func opAddSlice(
	p *pointerstructure.Pointer,
	parentVal reflect.Value,
	op *Operation,
	v interface{}) (interface{}, error) {
	// Get the final part. If the part is "-" we can directly set since
	// pointerstructure will handle the append.
	endPart := p.Parts[len(p.Parts)-1]
	if endPart == "-" {
		return p.Set(v, op.Value)
	}

	// "An element to add to an existing array - whereupon the supplied
	// value is added to the array at the indicated location.  Any
	// elements at or above the specified index are shifted one position
	// to the right."
	//
	// The above isn't a natural or built-in JSON pointer operation so
	// we're going to have to do some custom reflection here to mimic this:
	//
	// s = append(s, 0)
	// copy(s[i+1:], s[i:])
	// s[i] = x

	// First step: convert the part to an int so we can determine what index
	idxRaw, err := strconv.ParseInt(endPart, 10, 0)
	if err != nil {
		return v, fmt.Errorf("error parsing index %q: %s", endPart, err)
	}
	idx := int(idxRaw)

	// "The specified index MUST NOT be greater than the
	// number of elements in the array"
	if idx >= parentVal.Len() {
		return v, fmt.Errorf(
			"index %d is greater than the length %d",
			idx, parentVal.Len())
	}

	// Create a zero value to append for: s = append(s, 0)
	sliceType := parentVal.Type()
	slice := reflect.Append(parentVal, reflect.Indirect(reflect.New(sliceType.Elem())))

	// Perform the copy: copy(s[i+1:], s[i:])
	reflect.Copy(
		slice.Slice(idx+1, slice.Len()),
		slice.Slice(idx, slice.Len()))

	// Set the parent so that the slice is overwritten
	v, err = p.Parent().Set(v, slice.Interface())
	if err != nil {
		return v, err
	}

	// Write: s[i] = x
	return p.Set(v, op.Value)
}

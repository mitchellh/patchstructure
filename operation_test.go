package patchstructure

import (
	"fmt"
	"reflect"
	"testing"
)

func TestOperationApply(t *testing.T) {
	cases := []struct {
		Name      string
		Operation Operation
		Input     interface{}
		Expected  interface{}
		Err       bool
	}{
		//-----------------------------------------------------------
		// add
		//-----------------------------------------------------------

		// "The root of the target document - whereupon the specified value
		//  becomes the entire content of the target document."
		{
			"add: root",
			Operation{
				Op:    OpAdd,
				Path:  "",
				Value: "bar",
			},
			nil,
			"bar",
			false,
		},

		// "A member to add to an existing object - whereupon the supplied
		// value is added to that object at the indicated location.  If the
		// member already exists, it is replaced by the specified value."
		{
			"add: new member",
			Operation{
				Op:    OpAdd,
				Path:  "/a",
				Value: "bar",
			},
			map[string]interface{}{},
			map[string]interface{}{"a": "bar"},
			false,
		},

		{
			"add: existing member",
			Operation{
				Op:    OpAdd,
				Path:  "/a",
				Value: "bar",
			},
			map[string]interface{}{"a": "foo"},
			map[string]interface{}{"a": "bar"},
			false,
		},

		// "If the target location specifies an array index, a new value is
		// inserted into the array at the specified index."
		{
			"add: slice append",
			Operation{
				Op:    OpAdd,
				Path:  "/-",
				Value: "bar",
			},
			[]interface{}{1, 2},
			[]interface{}{1, 2, "bar"},
			false,
		},

		{
			"add: slice index",
			Operation{
				Op:    OpAdd,
				Path:  "/1",
				Value: "bar",
			},
			[]interface{}{1, 2},
			[]interface{}{1, "bar", 2},
			false,
		},

		{
			"add: slice index at 0",
			Operation{
				Op:    OpAdd,
				Path:  "/0",
				Value: "bar",
			},
			[]interface{}{1, 2},
			[]interface{}{"bar", 1, 2},
			false,
		},

		// "The specified index MUST NOT be greater than the
		// number of elements in the array"
		{
			"add: slice index out of bounds",
			Operation{
				Op:    OpAdd,
				Path:  "/4",
				Value: "bar",
			},
			[]interface{}{1, 2},
			nil,
			true,
		},

		// "However, the object itself or an array containing it does need to
		// exist, and it remains an error for that not to be the case."
		{
			"add: non-existent container",
			Operation{
				Op:    OpAdd,
				Path:  "/b/a",
				Value: "bar",
			},
			map[string]interface{}{"a": "foo"},
			nil,
			true,
		},

		//-----------------------------------------------------------
		// remove
		//-----------------------------------------------------------

		{
			"remove: map element",
			Operation{
				Op:   OpRemove,
				Path: "/foo/bar",
			},
			map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": 42,
				},
			},
			map[string]interface{}{
				"foo": map[string]interface{}{},
			},
			false,
		},

		{
			"remove: map element that doesn't exist",
			Operation{
				Op:   OpRemove,
				Path: "/foo/baz",
			},
			map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": 42,
				},
			},
			nil,
			true,
		},

		{
			"remove: slice index at 0",
			Operation{
				Op:   OpRemove,
				Path: "/0",
			},
			[]interface{}{1, 2},
			[]interface{}{2},
			false,
		},

		{
			"remove: slice index that doesn't exist",
			Operation{
				Op:   OpRemove,
				Path: "/4",
			},
			[]interface{}{1, 2},
			nil,
			true,
		},

		//-----------------------------------------------------------
		// replace
		//-----------------------------------------------------------

		{
			"replace: root",
			Operation{
				Op:    OpReplace,
				Path:  "",
				Value: "bar",
			},
			nil,
			"bar",
			false,
		},

		{
			"replace: new member",
			Operation{
				Op:    OpReplace,
				Path:  "/a",
				Value: "bar",
			},
			map[string]interface{}{},
			nil,
			true,
		},

		{
			"replace: existing member",
			Operation{
				Op:    OpReplace,
				Path:  "/a",
				Value: "bar",
			},
			map[string]interface{}{"a": "foo"},
			map[string]interface{}{"a": "bar"},
			false,
		},

		// NOTE(mitchellh): It is unclear what the RFC expects for this
		// behavior. It says that the target path must exist, and yet
		// I'm unsure if a "-" addr exists... it isn't clear.
		{
			"replace: slice append",
			Operation{
				Op:    OpReplace,
				Path:  "/-",
				Value: "bar",
			},
			[]interface{}{1, 2},
			nil,
			true,
		},

		{
			"replace: slice index",
			Operation{
				Op:    OpReplace,
				Path:  "/1",
				Value: "bar",
			},
			[]interface{}{1, 2},
			[]interface{}{1, "bar"},
			false,
		},

		{
			"replace: slice index at 0",
			Operation{
				Op:    OpReplace,
				Path:  "/0",
				Value: "bar",
			},
			[]interface{}{1, 2},
			[]interface{}{"bar", 2},
			false,
		},

		{
			"replace: slice index out of bounds",
			Operation{
				Op:    OpReplace,
				Path:  "/4",
				Value: "bar",
			},
			[]interface{}{1, 2},
			nil,
			true,
		},

		{
			"replace: non-existent container",
			Operation{
				Op:    OpReplace,
				Path:  "/b/a",
				Value: "bar",
			},
			map[string]interface{}{"a": "foo"},
			nil,
			true,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			actual, err := tc.Operation.Apply(tc.Input)
			if (err != nil) != tc.Err {
				t.Fatalf("err: %s", err)
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(actual, tc.Expected) {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}

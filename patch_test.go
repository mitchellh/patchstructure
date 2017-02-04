package patchstructure

import (
	"fmt"
	"reflect"
	"testing"
)

// Note that most operation tests that are more exhaustive are in
// operation_test.go. This test just tests basic sequences to ensure
// Patch behavior.
func TestPatch(t *testing.T) {
	cases := []struct {
		Name     string
		Ops      []*Operation
		Input    interface{}
		Expected interface{}
		Err      bool
	}{
		{
			"basic sequence",
			[]*Operation{
				&Operation{
					Op:    OpAdd,
					Path:  "/a",
					Value: "A",
				},

				&Operation{
					Op:   OpRemove,
					Path: "/b",
				},
			},
			map[string]interface{}{"b": 42},
			map[string]interface{}{"a": "A"},
			false,
		},

		{
			"partial failure",
			[]*Operation{
				&Operation{
					Op:    OpAdd,
					Path:  "/a",
					Value: "A",
				},

				&Operation{
					Op:   OpRemove,
					Path: "/c",
				},

				&Operation{
					Op:   OpRemove,
					Path: "/b",
				},
			},
			map[string]interface{}{"b": 42},
			map[string]interface{}{"a": "A", "b": 42},
			true,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			actual, err := Patch(tc.Input, tc.Ops)
			if (err != nil) != tc.Err {
				t.Fatalf("err: %s", err)
			}

			if !reflect.DeepEqual(actual, tc.Expected) {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}

package patchstructure

import (
	"fmt"
)

func ExamplePatch() {
	complex := map[string]interface{}{
		"alice": 42,
		"bob": []interface{}{
			map[string]interface{}{
				"name": "Bob",
			},
		},
	}

	value, err := Patch(complex, []*Operation{
		&Operation{
			Op:   OpCopy,
			Path: "/alice",
			From: "/bob",
		},

		&Operation{
			Op:    OpReplace,
			Path:  "/alice/0/name",
			Value: "Alice",
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", value)
	// Output:
	// map[alice:[map[name:Alice]] bob:[map[name:Bob]]]
}

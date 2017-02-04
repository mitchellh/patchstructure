# patchstructure [![GoDoc](https://godoc.org/github.com/mitchellh/patchstructure?status.svg)](https://godoc.org/github.com/mitchellh/patchstructure)

patchstructure is a Go library for applying patches to modify existing
Go structures.

patchstructure is based on
[JSON Patch (RFC 6902)](https://tools.ietf.org/html/rfc6902), but
applies to Go strucutures instead of JSON objects.

The goal of patchstructure is to provide a single API and format for
representing and applying changes to Go structures. With this in place,
diffs between structures can be represented, changes to a structure
can be treated as a log, etc.

## Features

  * Apply a "patch" to perform a set of operations on a Go structure

  * Operations support add, remove, replace, move, copy

  * Operations work on all Go primitive types and collection types

  * JSON encode/decode Operation structures

For an exhaustive list of supported features, please view the
[JSON Patch RFC (RFC 6902)](https://tools.ietf.org/html/rfc6902) which
this implements completely, but for Go structures. Exceptions to the RFC
are documented below.

## Differences from RFC 6902

RFC 6902 was created for JSON structures. Due to minor differences in
Go structures, there are some differences. These don't change the
function of the RFC, but are important to know when using this library:

  * The "copy" operation will perform a deep copy by default using
    [copystructure](https://github.com/mitchellh/copystructure). You can
    set the `Shallow` field to true on the operation to avoid this behavior.

  * The "test" operation is currently implemented with `reflect.DeepEqual`
    which is _not_ exactly correct according to the RFC. This will work fine
    for most basic primitives but the limitations of that approach should be
    known. In the future we should rework this.

## Installation

Standard `go get`:

```
$ go get github.com/mitchellh/patchstructure
```

## Usage & Example

For usage and examples see the [Godoc](http://godoc.org/github.com/mitchellh/patchstructure).

A quick code example is shown below:

```go
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
```

// Package patchstructure is a Go library for applying patches to modify
// existing Go structures.
//
// patchstructure is based on
// [JSON Patch (RFC 6902)](https://tools.ietf.org/html/rfc6902), but
// applies to Go strucutures instead of JSON objects.
//
// The goal of patchstructure is to provide a single API and format for
// representing and applying changes to Go structures. With this in place,
// diffs between structures can be represented, changes to a structure
// can be treated as a log, etc.
package patchstructure

// Patch applies the set of operations sequentially to the value v.
//
// Patch will halt at the first error. In this case, the returned value
// may be a partial value. This differs from the JSON Patch RFC which states
// that a patch should be atomic. Due to the complexity and cost in deep
// copying and the ability for the interface to store unsupported types
// such as functions (as long as they're not addressed it is okay), we defer
// this functionality to the end user.
//
// If you wish to deep copy your structures take a look at the "copystruture"
// library and call that prior to this.
func Patch(v interface{}, ops []*Operation) (result interface{}, err error) {
	result = v
	for _, op := range ops {
		result, err = op.Apply(result)
		if err != nil {
			return
		}
	}

	return
}

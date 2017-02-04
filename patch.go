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

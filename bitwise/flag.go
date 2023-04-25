package bitwise

// package bitwise - contains types with methods to make some bitwise operations

// Flag - is a type that contains a set of boolean flags.
// Each bit of Flag present a boolean flag.
type Flag uint

// This is example - how to setup boolean flag as bit order
const (
	// uint8
	FLAG_EXAMPLE_1 = 1 << uint(iota)
	FLAG_EXAMPLE_2
	FLAG_EXAMPLE_3
	FLAG_EXAMPLE_4
	FLAG_EXAMPLE_5
	FLAG_EXAMPLE_6
	FLAG_EXAMPLE_7
	FLAG_EXAMPLE_8
	// uint16
	FLAG_EXAMPLE_9
	FLAG_EXAMPLE_10
	FLAG_EXAMPLE_11
	FLAG_EXAMPLE_12
	FLAG_EXAMPLE_13
	FLAG_EXAMPLE_14
	FLAG_EXAMPLE_15
	FLAG_EXAMPLE_16
	// uint32
	// ...
	// uint64
	// ...
)

// Set a boolean flag in `on` state
func (f *Flag) Set(flag uint, on bool) {
	if on {
		*f = Flag(uint(*f) | flag)
		return
	}

	// off - bitclear
	*f = Flag(uint(*f) &^ flag)
}

// Has check boolean flag and return it state
func (f Flag) Has(flag uint) bool {
	return uint(f)&flag == flag
}

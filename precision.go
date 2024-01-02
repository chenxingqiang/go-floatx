package gofloatx

type Precision int

const (

	// PrecisionExact is for non-subnormals that don't drop bits during conversion.
	// All of these can round-trip.  Should always convert to float16.
	PrecisionExact Precision = iota

	// PrecisionUnknown is for subnormals that don't drop bits during conversion but
	// not all of these can round-trip so precision is unknown without more effort.
	// Only 2046 of these can round-trip and the rest cannot round-trip.
	PrecisionUnknown

	// PrecisionInexact is for dropped significand bits and cannot round-trip.
	// Some of these are subnormals. Cannot round-trip float32->float16->float32.
	PrecisionInexact

	// PrecisionUnderflow is for Underflows. Cannot round-trip float32->float16->float32.
	PrecisionUnderflow

	// PrecisionOverflow is for Overflows. Cannot round-trip float32->float16->float32.
	PrecisionOverflow
)

// Copyright 2019 Montgomery Edwards⁴⁴⁸ and Faye Amacker
//
// Special thanks to Kathryn Long for her Rust implementation
// of bfloat16 at github.com/starkat99/half-rs (MIT license)

// package gofloatx defines support for half-precision floating-point numbers.
package gofloatx

import (
	"math"
	"strconv"
)

// bfloat16 represents IEEE 754 half-precision floating-point numbers (binary16).
type bfloat16 uint16

// Precision indicates whether the conversion to bfloat16 is
// exact, subnormal without dropped bits, inexact, underflow, or overflow.

// ErrInvalidNaNValue indicates a NaN was not received.

type bfloat16Error string

func (e bfloat16Error) Error() string { return string(e) }

// FromNaN32ps converts nan to IEEE binary16 NaN while preserving both
// signaling and payload. Unlike Fromfloat32(), which can only return
// qNaN because it sets quiet bit = 1, this can return both sNaN and qNaN.
// If the result is infinity (sNaN with empty payload), then the
// lowest bit of payload is set to make the result a NaN.
// Returns ErrInvalidNaNValue and 0x7c01 (sNaN) if nan isn't IEEE 754 NaN.
// This function was kept simple to be able to inline.
func BF16FromNaN32ps(nan float32) (bfloat16, error) {
	const SNAN = bfloat16(uint16(0x7c01)) // signaling NaN

	u32 := math.Float32bits(nan)
	sign := u32 & 0x80000000
	exp := u32 & 0x7f800000
	coef := u32 & 0x007fffff

	if (exp != 0x7f800000) || (coef == 0) {
		return SNAN, ErrInvalidNaNValue
	}

	u16 := uint16((sign >> 16) | uint32(0x7c00) | (coef >> 13))

	if (u16 & 0x03ff) == 0 {
		// result became infinity, make it NaN by setting lowest bit in payload
		u16 |= 0x0001
	}

	return bfloat16(u16), nil
}

// Float32 returns a float32 converted from f (bfloat16).
// This is a lossless conversion.
func (f bfloat16) Float32() float32 {
	u32 := f16bitsToF32bits(uint16(f))
	return math.Float32frombits(u32)
}

// Bits returns the IEEE 754 binary16 representation of f, with the sign bit
// of f and the result in the same bit position. Bits(Frombits(x)) == x.
func (f bfloat16) Bits() uint16 {
	return uint16(f)
}

// IsNaN reports whether f is an IEEE 754 binary16 “not-a-number” value.
func (f bfloat16) IsNaN() bool {
	return (f&0x7c00 == 0x7c00) && (f&0x03ff != 0)
}

// IsQuietNaN reports whether f is a quiet (non-signaling) IEEE 754 binary16
// “not-a-number” value.
func (f bfloat16) IsQuietNaN() bool {
	return (f&0x7c00 == 0x7c00) && (f&0x03ff != 0) && (f&0x0200 != 0)
}

// IsInf reports whether f is an infinity (inf).
// A sign > 0 reports whether f is positive inf.
// A sign < 0 reports whether f is negative inf.
// A sign == 0 reports whether f is either inf.
func (f bfloat16) IsInf(sign int) bool {
	return ((f == 0x7c00) && sign >= 0) ||
		(f == 0xfc00 && sign <= 0)
}

// IsFinite returns true if f is neither infinite nor NaN.
func (f bfloat16) IsFinite() bool {
	return (uint16(f) & uint16(0x7c00)) != uint16(0x7c00)
}

// IsNormal returns true if f is neither zero, infinite, subnormal, or NaN.
func (f bfloat16) IsNormal() bool {
	exp := uint16(f) & uint16(0x7c00)
	return (exp != uint16(0x7c00)) && (exp != 0)
}

// Signbit reports whether f is negative or negative zero.
func (f bfloat16) Signbit() bool {
	return (uint16(f) & uint16(0x8000)) != 0
}

// String satisfies the fmt.Stringer interface.
func (f bfloat16) String() string {
	return strconv.FormatFloat(float64(f.Float32()), 'f', -1, 32)
}

// f16bitsToF32bits returns uint32 (float32 bits) converted from specified uint16.
func bf16bitsToF32bits(in uint16) uint32 {
	// All 65536 conversions with this were confirmed to be correct
	// by Montgomery Edwards⁴⁴⁸ (github.com/x448).

	sign := uint32(in&0x8000) << 16 // sign for 32-bit
	exp := uint32(in&0x7c00) >> 10  // exponenent for 16-bit
	coef := uint32(in&0x03ff) << 13 // significand for 32-bit

	if exp == 0x1f {
		if coef == 0 {
			// infinity
			return sign | 0x7f800000 | coef
		}
		// NaN
		return sign | 0x7fc00000 | coef
	}

	if exp == 0 {
		if coef == 0 {
			// zero
			return sign
		}

		// normalize subnormal numbers
		exp++
		for coef&0x7f800000 == 0 {
			coef <<= 1
			exp--
		}
		coef &= 0x007fffff
	}

	return sign | ((exp + (0x7f - 0xf)) << 23) | coef
}

// f32bitsToF16bits returns uint16 (bfloat16 bits) converted from the specified float32.
// Conversion rounds to nearest integer with ties to even.
func f32bitsToBF16bits(u32 uint32) uint16 {
	// Translated from Rust to Go by Montgomery Edwards⁴⁴⁸ (github.com/x448).
	// All 4294967296 conversions with this were confirmed to be correct by x448.
	// Original Rust implementation is by Kathryn Long (github.com/starkat99) with MIT license.

	sign := u32 & 0x80000000
	exp := u32 & 0x7f800000
	coef := u32 & 0x007fffff

	if exp == 0x7f800000 {
		// NaN or Infinity
		nanBit := uint32(0)
		if coef != 0 {
			nanBit = uint32(0x0200)
		}
		return uint16((sign >> 16) | uint32(0x7c00) | nanBit | (coef >> 13))
	}

	halfSign := sign >> 16

	unbiasedExp := int32(exp>>23) - 127
	halfExp := unbiasedExp + 15

	if halfExp >= 0x1f {
		return uint16(halfSign | uint32(0x7c00))
	}

	if halfExp <= 0 {
		if 14-halfExp > 24 {
			return uint16(halfSign)
		}
		c := coef | uint32(0x00800000)
		halfCoef := c >> uint32(14-halfExp)
		roundBit := uint32(1) << uint32(13-halfExp)
		if (c&roundBit) != 0 && (c&(3*roundBit-1)) != 0 {
			halfCoef++
		}
		return uint16(halfSign | halfCoef)
	}

	uHalfExp := uint32(halfExp) << 10
	halfCoef := coef >> 13
	roundBit := uint32(0x00001000)
	if (coef&roundBit) != 0 && (coef&(3*roundBit-1)) != 0 {
		return uint16((halfSign | uHalfExp | halfCoef) + 1)
	}
	return uint16(halfSign | uHalfExp | halfCoef)
}

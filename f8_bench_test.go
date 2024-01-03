package floatx_test

import (
	floatx "github.com/chenxingqiang/go-floatx"
	"math"
	"testing"
)

// prevent comPiler optimizing out code by assigning to these
var F8ResultF8 floatx.Float8
var F8ResultF32 float32
var F8ResultStr string
var F8PCN floatx.F8Precision

func F8BenchmarkFloat32Pi(b *testing.B) {
	result := float32(0)
	Pi32 := float32(math.Pi)
	Pi8 := floatx.F8Fromfloat32(Pi32)
	for i := 0; i < b.N; i++ {
		f8 := floatx.F8Frombits(uint8(Pi8))
		result = f8.Float32()
	}
	F8ResultF32 = result
}

func F8BenchmarkFrombits(b *testing.B) {
	result := floatx.Float8(0)
	Pi32 := float32(math.Pi)
	Pi8 := floatx.F8Fromfloat32(Pi32)
	for i := 0; i < b.N; i++ {
		result = floatx.F8Frombits(uint8(Pi8))
	}
	F8ResultF8 = result
}

func F8BenchmarkFromFloat32Pi(b *testing.B) {
	result := floatx.Float8(0)

	Pi := float32(math.Pi)
	for i := 0; i < b.N; i++ {
		result = floatx.F8Fromfloat32(Pi)
	}
	F8ResultF8 = result
}

func F8BenchmarkFromFloat32nan(b *testing.B) {
	result := floatx.Float8(0)

	nan := float32(math.NaN())
	for i := 0; i < b.N; i++ {
		result = floatx.F8Fromfloat32(nan)
	}
	F8ResultF8 = result
}

func F8BenchmarkFromFloat32subnorm(b *testing.B) {
	result := floatx.Float8(0)

	subnorm := math.Float32frombits(0x007fffff)
	for i := 0; i < b.N; i++ {
		result = floatx.F8Fromfloat32(subnorm)
	}
	F8ResultF8 = result
}

func F8BenchmarkPrecisionFromFloat32(b *testing.B) {
	var result floatx.F8Precision

	for i := 0; i < b.N; i++ {
		f32 := float32(0.00001) + float32(0.00001)
		result = floatx.F8PrecisionFromfloat32(f32)
	}
	F8PCN = result
}

func F8BenchmarkString(b *testing.B) {
	var result string

	Pi32 := float32(math.Pi)
	Pi8 := floatx.F8Fromfloat32(Pi32)
	for i := 0; i < b.N; i++ {
		result = Pi8.String()
	}
	F8ResultStr = result
}

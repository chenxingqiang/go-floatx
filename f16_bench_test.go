package floatx_test

import (
	floatx "github.com/chenxingqiang/go-floatx"
	"math"
	"testing"
)

// prevent comPiler optimizing out code by assigning to these
var F16ResultF16 floatx.Float16
var F16ResultF32 float32
var F16ResultStr string
var F16PCN floatx.F16Precision

func F16BenchmarkFloat32Pi(b *testing.B) {
	result := float32(0)
	Pi32 := float32(math.Pi)
	Pi16 := floatx.F16Fromfloat32(Pi32)
	for i := 0; i < b.N; i++ {
		f16 := floatx.F16Frombits(uint16(Pi16))
		result = f16.Float32()
	}
	F16ResultF32 = result
}

func F16BenchmarkFrombits(b *testing.B) {
	result := floatx.Float16(0)
	Pi32 := float32(math.Pi)
	Pi16 := floatx.F16Fromfloat32(Pi32)
	for i := 0; i < b.N; i++ {
		result = floatx.F16Frombits(uint16(Pi16))
	}
	F16ResultF16 = result
}

func F16BenchmarkFromFloat32Pi(b *testing.B) {
	result := floatx.Float16(0)

	Pi := float32(math.Pi)
	for i := 0; i < b.N; i++ {
		result = floatx.F16Fromfloat32(Pi)
	}
	F16ResultF16 = result
}

func F16BenchmarkFromFloat32nan(b *testing.B) {
	result := floatx.Float16(0)

	nan := float32(math.NaN())
	for i := 0; i < b.N; i++ {
		result = floatx.F16Fromfloat32(nan)
	}
	F16ResultF16 = result
}

func F16BenchmarkFromFloat32subnorm(b *testing.B) {
	result := floatx.Float16(0)

	subnorm := math.Float32frombits(0x007fffff)
	for i := 0; i < b.N; i++ {
		result = floatx.F16Fromfloat32(subnorm)
	}
	F16ResultF16 = result
}

func F16BenchmarkPrecisionFromFloat32(b *testing.B) {
	var result floatx.F16Precision

	for i := 0; i < b.N; i++ {
		f32 := float32(0.00001) + float32(0.00001)
		result = floatx.F16PrecisionFromfloat32(f32)
	}
	F16PCN = result
}

func F16BenchmarkString(b *testing.B) {
	var result string

	Pi32 := float32(math.Pi)
	Pi16 := floatx.F16Fromfloat32(Pi32)
	for i := 0; i < b.N; i++ {
		result = Pi16.String()
	}
	F16ResultStr = result
}

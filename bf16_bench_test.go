package floatx_test

import (
	floatx "github.com/chenxingqiang/go-floatx"
	"math"
	"testing"
)

// prevent comPiler optimizing out code by assigning to these
var BF16ResultBF16 floatx.BFloat16
var BF16ResultF32 float32
var BF16ResultStr string
var BF16PCN floatx.BF16Precision

func BF16BenchmarkFloat32Pi(b *testing.B) {
	result := float32(0)
	Pi32 := float32(math.Pi)
	Pi16 := floatx.BF16Fromfloat32(Pi32)
	for i := 0; i < b.N; i++ {
		f16 := floatx.BF16Frombits(uint16(Pi16))
		result = f16.Float32()
	}
	resultF32 = result
}

func BF16BenchmarkFrombits(b *testing.B) {
	result := floatx.BFloat16(0)
	Pi32 := float32(math.Pi)
	Pi16 := floatx.BF16Fromfloat32(Pi32)
	for i := 0; i < b.N; i++ {
		result = floatx.BF16Frombits(uint16(Pi16))
	}
	BF16ResultBF16 = result
}

func BF16BenchmarkFromFloat32Pi(b *testing.B) {
	result := floatx.BFloat16(0)

	Pi := float32(math.Pi)
	for i := 0; i < b.N; i++ {
		result = floatx.BF16Fromfloat32(Pi)
	}
	BF16ResultBF16 = result
}

func BF16BenchmarkFromFloat32nan(b *testing.B) {
	result := floatx.Float16(0)

	nan := float32(math.NaN())
	for i := 0; i < b.N; i++ {
		result = floatx.F16Fromfloat32(nan)
	}
	resultF16 = result
}

func BF16BenchmarkFromFloat32subnorm(b *testing.B) {
	result := floatx.BFloat16(0)

	subnorm := math.Float32frombits(0x007fffff)
	for i := 0; i < b.N; i++ {
		result = floatx.BF16Fromfloat32(subnorm)
	}
	BF16ResultBF16 = result
}

func BF16BenchmarkPrecisionFromFloat32(b *testing.B) {
	var result floatx.BF16Precision

	for i := 0; i < b.N; i++ {
		f32 := float32(0.00001) + float32(0.00001)
		result = floatx.BF16PrecisionFromfloat32(f32)
	}
	BF16PCN = result
}

func BF16BenchmarkString(b *testing.B) {
	var result string

	Pi32 := float32(math.Pi)
	Pi16 := floatx.F16Fromfloat32(Pi32)
	for i := 0; i < b.N; i++ {
		result = Pi16.String()
	}
	BF16ResultStr = result
}

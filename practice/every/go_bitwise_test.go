package every

import "testing"

func TestBitOperations(t *testing.T) {

	tests := []struct {
		name      string
		a         uint8
		b         uint8
		operation func(uint8, uint8) uint8
		expected  uint8
	}{
		{
			name:      "AND",
			a:         12, // 1100
			b:         10, // 1010
			operation: BitAnd,
			expected:  8, // 1000
		},
		{
			name:      "OR",
			a:         12, // 1100
			b:         10, // 1010
			operation: BitOr,
			expected:  14, // 1110
		},
		{
			name:      "XOR",
			a:         12, // 1100
			b:         10, // 1010
			operation: BitXor,
			expected:  6, // 0110
		},
		{
			name:      "AND NOT",
			a:         12, // 1100
			b:         10, // 1010
			operation: BitNot,
			expected:  4, // 0100
		},
		{

			name:      "AND with zero",
			a:         0,  // 00000000
			b:         10, // 00001010
			operation: BitAnd,
			expected:  0, // 00000000
		},
		{
			name:      "OR with zero",
			a:         0,  // 00000000
			b:         10, // 00001010
			operation: BitOr,
			expected:  10, // 00001010
		},
		{
			name:      "XOR with zero",
			a:         0,  // 00000000
			b:         10, // 00001010
			operation: BitXor,
			expected:  10, // 00001010
		},
		{
			name:      "AND NOT with zero",
			a:         0,  // 00000000
			b:         10, // 00001010
			operation: BitNot,
			expected:  0, // 00000000
		},
		{
			name:      "AND with full byte",
			a:         255, // 11111111
			b:         10,  // 00001010
			operation: BitAnd,
			expected:  10, // 00001010
		},
		{
			name:      "OR with full byte",
			a:         255, // 11111111
			b:         10,  // 00001010
			operation: BitOr,
			expected:  255, // 11111111
		},
		{
			name:      "XOR with full byte",
			a:         255, // 11111111
			b:         10,  // 00001010
			operation: BitXor,
			expected:  245, // 11110101
		},
		{
			name:      "AND NOT with full byte",
			a:         255, // 11111111
			b:         10,  // 00001010
			operation: BitNot,
			expected:  245, // 11110101
		},
	}
	for _, test := range tests {
		result := test.operation(test.a, test.b)
		if result != test.expected {
			t.Errorf(
				"%s: got %d / %08b, want %d / %08b",
				test.name,
				result,
				result,
				test.expected,
				test.expected,
			)
		}
	}

}

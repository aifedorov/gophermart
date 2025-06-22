package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidOrderNumber(t *testing.T) {
	t.Parallel()

	type want struct {
		result bool
	}
	tests := []struct {
		name   string
		number string
		want   want
	}{
		{
			name:   "valid order number 1",
			number: "4532015112830366",
			want: want{
				result: true,
			},
		},
		{
			name:   "valid order number 2",
			number: "5555555555554444",
			want: want{
				result: true,
			},
		},
		{
			name:   "valid order number with spaces",
			number: "4532 0151 1283 0366",
			want: want{
				result: true,
			},
		},
		{
			name:   "invalid order number",
			number: "4532015112830367",
			want: want{
				result: false,
			},
		},
		{
			name:   "empty order number",
			number: "",
			want: want{
				result: false,
			},
		},
		{
			name:   "contains letters",
			number: "453201511283036a",
			want: want{
				result: false,
			},
		},
		{
			name:   "contains special characters",
			number: "4532-0151-1283-0366",
			want: want{
				result: false,
			},
		},
		{
			name:   "single digit valid",
			number: "0",
			want: want{
				result: true,
			},
		},
		{
			name:   "single digit invalid",
			number: "1",
			want: want{
				result: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := isValidOrderNumber(tt.number)
			assert.Equal(t, tt.want.result, result)
		})
	}
}

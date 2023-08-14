package convbytes

import (
	"math"
	"testing"
)

func TestConvert(t *testing.T) {
	type args struct {
		value    int64
		fromUnit int64
		toUnit   int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "MB to GB",
			args: args{
				value:    1024,
				fromUnit: MB,
				toUnit:   GB,
			},
			want: 1,
		},
		{
			name: "KB to GB",
			args: args{
				value:    1024 * 1024,
				fromUnit: KB,
				toUnit:   GB,
			},
			want: 1,
		},
		{
			name: "GB to KB",
			args: args{
				value:    1,
				fromUnit: GB,
				toUnit:   KB,
			},
			want: 1024 * 1024,
		},
		{
			name: "TB to GB",
			args: args{
				value:    1,
				fromUnit: TB,
				toUnit:   GB,
			},
			want: 1024,
		},
		{
			name: "Negative value",
			args: args{
				value:    -1,
				fromUnit: GB,
				toUnit:   KB,
			},
			want: -1024 * 1024,
		},
		{
			name: "Big int MB to GB",
			args: args{
				value:    90368775807,
				fromUnit: MB,
				toUnit:   GB,
			},
			want: 90368775807 / 1024,
		},
		{
			name: "PB to TB",
			args: args{
				value:    1,
				fromUnit: PB,
				toUnit:   TB,
			},
			want: 1024,
		},
		{
			name: "PB to bytes",
			args: args{
				value:    16,
				fromUnit: PB,
				toUnit:   B,
			},
			want: int64(16 * math.Pow(1024, 5)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Convert(tt.args.value, tt.args.fromUnit, tt.args.toUnit); got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromBytes(t *testing.T) {
	type args struct {
		bytes int64
		unit  int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Bytes to MB",
			args: args{
				bytes: 1024 * 1024,
				unit:  MB,
			},
			want: 1,
		},
		{
			name: "Bytes to GB",
			args: args{
				bytes: 1024 * 1024 * 1024,
				unit:  GB,
			},
			want: 1,
		},
		{
			name: "Bytes to TB",
			args: args{
				bytes: 1024 * 1024 * 1024 * 1024,
				unit:  TB,
			},
			want: 1,
		},
		{
			name: "Bytes to B",
			args: args{
				bytes: 500,
				unit:  B,
			},
			want: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromBytes(tt.args.bytes, tt.args.unit); got != tt.want {
				t.Errorf("FromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "String without unit",
			args: args{
				s: "1024",
			},
			want: 1024,
		},
		{
			name: "String with MB",
			args: args{
				s: "1 MB",
			},
			want: 1024 * 1024,
		},
		{
			name: "String with GB",
			args: args{
				s: "1 GB",
			},
			want: 1024 * 1024 * 1024,
		},
		{
			name: "String with whitespace",
			args: args{
				s: "  1   GB   ",
			},
			want: 1024 * 1024 * 1024,
		},
		{
			name: "String with TB",
			args: args{
				s: "1TB",
			},
			want: 1024 * 1024 * 1024 * 1024,
		},
		{
			name: "String with lowercase unit",
			args: args{
				s: "1 gb",
			},
			want: 1024 * 1024 * 1024,
		},
		{
			name: "Valid string with PB",
			args: args{s: "1 PB"},
			want: 1024 * 1024 * 1024 * 1024 * 1024,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrToBytes(tt.args.s); got != tt.want {
				t.Errorf("StrToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromStringOk(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 bool
	}{
		{
			name:  "Valid string without unit",
			args:  args{s: "1024"},
			want:  1024,
			want1: true,
		},
		{
			name:  "Valid string with MB",
			args:  args{s: "1 MB"},
			want:  1024 * 1024,
			want1: true,
		},
		{
			name:  "Invalid string format",
			args:  args{s: "invalid"},
			want:  0,
			want1: false,
		},
		{
			name:  "String with invalid unit",
			args:  args{s: "1 ZB"},
			want:  0,
			want1: false,
		},
		{
			name:  "String with leading and trailing spaces",
			args:  args{s: "   1024 KB   "},
			want:  1024 * 1024,
			want1: true,
		},
		{
			name:  "Empty string",
			args:  args{s: ""},
			want:  0,
			want1: false,
		},
		{
			name:  "Valid string with PB",
			args:  args{s: "1 PB"},
			want:  1024 * 1024 * 1024 * 1024 * 1024,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := StrToBytesOk(tt.args.s)
			if got != tt.want {
				t.Errorf("StrToBytesOk() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("StrToBytesOk() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestToBytes(t *testing.T) {
	type args struct {
		value int64
		unit  int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "MB to bytes",
			args: args{
				value: 1,
				unit:  MB,
			},
			want: 1024 * 1024,
		},
		{
			name: "GB to bytes",
			args: args{
				value: 1,
				unit:  GB,
			},
			want: 1024 * 1024 * 1024,
		},
		{
			name: "TB to bytes",
			args: args{
				value: 1,
				unit:  TB,
			},
			want: 1024 * 1024 * 1024 * 1024,
		},
		{
			name: "B to bytes",
			args: args{
				value: 500,
				unit:  B,
			},
			want: 500,
		},
		{
			name: "Negative value",
			args: args{
				value: -1,
				unit:  GB,
			},
			want: -1024 * 1024 * 1024,
		},
		{
			name: "PB to bytes",
			args: args{
				value: 1,
				unit:  PB,
			},
			want: 1024 * 1024 * 1024 * 1024 * 1024,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBytes(tt.args.value, tt.args.unit); got != tt.want {
				t.Errorf("ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToUnit(t *testing.T) {
	type args struct {
		s          string
		targetUnit int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "String without unit to MB",
			args: args{
				s:          "1024",
				targetUnit: MB,
			},
			want: 1024,
		},
		{
			name: "String with MB to GB",
			args: args{
				s:          "1024 MB",
				targetUnit: GB,
			},
			want: 1,
		},
		{
			name: "String with GB to TB",
			args: args{
				s:          "1024 GB",
				targetUnit: TB,
			},
			want: 1,
		},
		{
			name: "String with whitespace to GB",
			args: args{
				s:          "  1024   GB   ",
				targetUnit: TB,
			},
			want: 1,
		},
		{
			name: "String without unit to GB",
			args: args{
				s:          "1024",
				targetUnit: GB,
			},
			want: 1024,
		},
		{
			name: "String with invalid format",
			args: args{
				s:          "invalid",
				targetUnit: GB,
			},
			want: 0, // or whatever your default/error return value is
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrToUnit(tt.args.s, tt.args.targetUnit); got != tt.want {
				t.Errorf("StrToUnit() = %v, want %v", got, tt.want)
			}
		})
	}
}

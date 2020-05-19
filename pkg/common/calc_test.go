package common

import (
	"testing"
)

func TestMin(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.nums...); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args struct {
		nums []int
	}
	vv := make(map[string]int)

	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "s001",
			args: args{
				[]int{1, 10, 9, 8},
			},
			want: 10,
		},
		{
			name: "s002",
			args: args{
				[]int{vv["aa"], vv["vv"], 9, 8},
			},
			want: 9,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.nums...); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtoi(t *testing.T) {
	type args struct {
		s            string
		defaultvalue int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Atoi(tt.args.s, tt.args.defaultvalue); got != tt.want {
				t.Errorf("Atoi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundI(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundI(tt.args.f); got != tt.want {
				t.Errorf("RoundI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRound(t *testing.T) {
	type args struct {
		f      float64
		places int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Round(tt.args.f, tt.args.places); got != tt.want {
				t.Errorf("Round() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCeil(t *testing.T) {
	type args struct {
		f      float64
		places int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ceil(tt.args.f, tt.args.places); got != tt.want {
				t.Errorf("Ceil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRound5sha(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "正常系01",
			args: args{
				f: 6.5,
			},
			want: 6.0,
		},
		{
			name: "正常系02",
			args: args{
				f: 6.6,
			},
			want: 7.0,
		},
		{
			name: "正常系03",
			args: args{
				f: 6.51,
			},
			want: 7.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Round5sha(tt.args.f); got != tt.want {
				t.Errorf("Round5sha() = %v, want %v", got, tt.want)
			}
		})
	}
}

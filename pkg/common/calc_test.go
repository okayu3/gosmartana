package common

import "testing"

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

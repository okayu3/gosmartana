package common

import "testing"

func TestYmdW2g(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "s001",
			args: args{
				s: "5010502",
			},
			want: "20190502",
		},
		{
			name: "s002",
			args: args{
				s: "3470403",
			},
			want: "19720403",
		},
		{
			name: "s003",
			args: args{
				s: "4301231",
			},
			want: "20181231",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := YmdW2g(tt.args.s); got != tt.want {
				t.Errorf("YmdW2g() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYmW2g(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := YmW2g(tt.args.s); got != tt.want {
				t.Errorf("YmW2g() = %v, want %v", got, tt.want)
			}
		})
	}
}

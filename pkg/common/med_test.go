package common

import "testing"

func Test_isDoubtDisease(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "s001",
			args: args{
				s: "999980029999",
			},
			want: true,
		},
		{
			name: "s002",
			args: args{
				s: "8002",
			},
			want: true,
		},
		{
			name: "s003",
			args: args{
				s: "",
			},
			want: false,
		},
		{
			name: "s004",
			args: args{
				s: "980029",
			},
			want: false,
		},
		{
			name: "s005",
			args: args{
				s: "000080029",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDoubtDisease(tt.args.s); got != tt.want {
				t.Errorf("isDoubtDisease() = %v, want %v", got, tt.want)
			}
		})
	}
}

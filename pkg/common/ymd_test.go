package common

import (
	"testing"
)

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

func TestAgeAt(t *testing.T) {
	type args struct {
		ymdB  string
		atYmd string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "正常系01",
			args: args{
				ymdB:  "19720403",
				atYmd: "20200404",
			},
			want: 48,
		},
		{
			name: "正常系02",
			args: args{
				ymdB:  "19720403",
				atYmd: "20200403",
			},
			want: 48,
		},
		{
			name: "正常系03",
			args: args{
				ymdB:  "19720403",
				atYmd: "20200402",
			},
			want: 47,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AgeAt(tt.args.ymdB, tt.args.atYmd); got != tt.want {
				t.Errorf("AgeAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnualAtYmd(t *testing.T) {
	type args struct {
		atYmd string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "正常系01",
			args: args{
				atYmd: "20200401",
			},
			want: 2020,
		},
		{
			name: "正常系02",
			args: args{
				atYmd: "20200402",
			},
			want: 2020,
		},
		{
			name: "正常系03",
			args: args{
				atYmd: "20200331",
			},
			want: 2019,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnnualAtYmd(tt.args.atYmd); got != tt.want {
				t.Errorf("AnnualAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnnualAtYm(t *testing.T) {
	type args struct {
		atYm string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "正常系01",
			args: args{
				atYm: "202012",
			},
			want: 2020,
		},
		{
			name: "正常系02",
			args: args{
				atYm: "202004",
			},
			want: 2020,
		},
		{
			name: "正常系03",
			args: args{
				atYm: "202003",
			},
			want: 2019,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnnualAtYm(tt.args.atYm); got != tt.want {
				t.Errorf("AnnualAtYm() = %v, want %v", got, tt.want)
			}
		})
	}
}

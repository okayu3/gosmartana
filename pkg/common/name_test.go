package common

import "testing"

func TestDevideName(t *testing.T) {
	type args struct {
		nm string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
		{
			name: "正常系１",
			args: args{
				nm: "前進　太郎",
			},
			want:  "前進",
			want1: "太郎",
		},
		{
			name: "正常系２",
			args: args{
				nm: "前進　パウロ　太郎",
			},
			want:  "前進",
			want1: "パウロ　太郎",
		},
		{
			name: "異常系１",
			args: args{
				nm: "前進太郎",
			},
			want:  "前進太郎",
			want1: Empty,
		},
		{
			name: "異常系２",
			args: args{
				nm: "　前進太郎",
			},
			want:  "前進太郎",
			want1: Empty,
		},
		{
			name: "異常系３",
			args: args{
				nm: "",
			},
			want:  Empty,
			want1: Empty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := DevideName(tt.args.nm)
			if got != tt.want {
				t.Errorf("DevideName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DevideName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

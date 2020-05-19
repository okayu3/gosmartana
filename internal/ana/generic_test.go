package ana

import (
	"reflect"
	"testing"
)

func setup() {
	fnmAnyDrg := "C:/task/prj/YG01/ref/cleansing/2020/etc/generic/A015_01_mst_any_yakka_period.csv"
	fnmGeneFlg := "C:/task/prj/YG01/ref/cleansing/2020/etc/generic/A015_01_mst_genestat_period.csv"
	fnmCheap := "C:/task/prj/YG01/ref/cleansing/2020/etc/generic/A015_01_mst_cheapest_period_202004.csv"
	fnmExpensv := "C:/task/prj/YG01/ref/cleansing/2020/etc/generic/A015_01_mst_expensive_period_202004.csv"
	LoadMstGeneric(fnmAnyDrg, fnmGeneFlg, fnmCheap, fnmExpensv)
}

func TestLoadMstGeneric(t *testing.T) {
	if len(mstDrgCheap) <= 0 {
		setup()
	}
	type args struct {
		fnmAnyDrg  string
		fnmGeneFlg string
		fnmCheap   string
		fnmExpensv string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadMstGeneric(tt.args.fnmAnyDrg, tt.args.fnmGeneFlg, tt.args.fnmCheap, tt.args.fnmExpensv)
		})
	}
}

func Test_yakkaAt(t *testing.T) {
	if len(mstDrgCheap) <= 0 {
		setup()
	}
	type args struct {
		drcd string
		ym   string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "正常系テスト01",
			args: args{
				drcd: "610406013",
				ym:   "201312",
			},
			want: 6.2,
		},
		{
			name: "正常系テスト02",
			args: args{
				drcd: "610406013",
				ym:   "201403",
			},
			want: 6.2,
		},
		{
			name: "正常系テスト03",
			args: args{
				drcd: "610406013",
				ym:   "201404",
			},
			want: 5.8,
		},
		{
			name: "正常系テスト04",
			args: args{
				drcd: "610406013",
				ym:   "202004",
			},
			want: 5.8,
		},
		{
			name: "異常系テスト01",
			args: args{
				drcd: "610406013",
				ym:   "212004",
			},
			want: 0,
		},
		{
			name: "異常系テスト02",
			args: args{
				drcd: "610406013",
				ym:   "192004",
			},
			want: 0,
		},
		{
			name: "異常系テスト03",
			args: args{
				drcd: "610406013",
				ym:   "198904",
			},
			want: 6.2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := yakkaAt(tt.args.drcd, tt.args.ym); got != tt.want {
				t.Errorf("yakkaAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_geneFlgAt(t *testing.T) {
	if len(mstDrgCheap) <= 0 {
		setup()
	}
	type args struct {
		drcd string
		ym   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "正常系テスト01",
			args: args{
				drcd: "620000199",
				ym:   "201905",
			},
			want: "3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := geneFlgAt(tt.args.drcd, tt.args.ym); got != tt.want {
				t.Errorf("geneFlgAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_geneCheapExpsvAt(t *testing.T) {
	if len(mstDrgCheap) <= 0 {
		setup()
	}
	type args struct {
		drcd string
		ym   string
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 float64
	}{
		// TODO: Add test cases.
		{
			name: "正常系テスト01",
			args: args{
				drcd: "620000199",
				ym:   "201905",
			},
			want:  487,
			want1: 576,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := geneCheapExpsvAt(tt.args.drcd, tt.args.ym)
			if got != tt.want {
				t.Errorf("geneCheapExpsvAt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("geneCheapExpsvAt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_newYakkaRange(t *testing.T) {
	type args struct {
		one string
	}
	tests := []struct {
		name string
		args args
		want yakkaRange
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newYakkaRange(tt.args.one); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newYakkaRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newGeneFlgRange(t *testing.T) {
	type args struct {
		one string
	}
	tests := []struct {
		name string
		args args
		want geneFlgRange
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newGeneFlgRange(tt.args.one); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newGeneFlgRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

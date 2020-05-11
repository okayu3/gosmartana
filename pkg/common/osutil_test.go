package common

import (
	"reflect"
	"testing"
)

func TestFileExists(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.args.filename); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MakeDir(tt.args.path)
		})
	}
}

func TestListUpFiles(t *testing.T) {
	type args struct {
		root   string
		prefix string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "正常系：",
			args: args{
				root:   "C:/task/prj/YG01/wk",
				prefix: "Q",
				suffix: ".pl",
			},
			want: []string{"C:/task/prj/YG01/wk/Q01_001_constrece_forGolang.pl",
				"C:/task/prj/YG01/wk/Q02_tosekimst_forGo.pl"},
		},
		{
			name: "正常系：prefix is empty",
			args: args{
				root:   "C:/task/prj/YG01/wk",
				prefix: "",
				suffix: ".pl",
			},
			want: []string{"C:/task/prj/YG01/wk/Q01_001_constrece_forGolang.pl",
				"C:/task/prj/YG01/wk/Q02_tosekimst_forGo.pl"},
		},
		{
			name: "正常系：suffix is empty",
			args: args{
				root:   "C:/task/prj/YG01/wk",
				prefix: "Q01",
				suffix: "",
			},
			want: []string{"C:/task/prj/YG01/wk/Q01_001_constrece_forGolang.pl",
				"C:/task/prj/YG01/wk/Q01_001_constrece_forGolang.pl.1~"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListUpFiles(tt.args.root, tt.args.prefix, tt.args.suffix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListUpFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListUpRece(t *testing.T) {
	type args struct {
		root   string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		// TODO: Add test cases.
		{
			name: "正常系：",
			args: args{
				root:   "C:/task/prj/YG01/sample/rece/H29.05月診療分/201706/06140248/0_COMMON001",
				suffix: "",
			},
			want: [][]string{
				{"C:/task/prj/YG01/sample/rece/H29.05月診療分/201706/06140248/0_COMMON001/11_RECODEINFO_MED.CSV"},
				{"C:/task/prj/YG01/sample/rece/H29.05月診療分/201706/06140248/0_COMMON001/12_RECODEINFO_DPC.CSV"},
				{"C:/task/prj/YG01/sample/rece/H29.05月診療分/201706/06140248/0_COMMON001/13_RECODEINFO_DEN.CSV"},
				{"C:/task/prj/YG01/sample/rece/H29.05月診療分/201706/06140248/0_COMMON001/14_RECODEINFO_PHA.CSV"},
				{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListUpRece(tt.args.root, tt.args.suffix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListUpRece() = %v, want %v", got, tt.want)
			}
		})
	}
}

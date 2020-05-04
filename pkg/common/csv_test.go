package common

import (
	"fmt"
	"os"
	"testing"
)

func Test_failOnError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			failOnError(tt.args.err)
		})
	}
}

func Test_loadCsvMain(t *testing.T) {
	type args struct {
		in   *os.File
		r    CsvCallback
		mode int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadCsvMain(tt.args.in, tt.args.r, tt.args.mode)
		})
	}
}

func TestLoadCSV(t *testing.T) {
	type args struct {
		fnm  string
		r    CsvCallback
		mode int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "CSVファイル",
			args: args{
				fnm: "C:/task/garden/py/smartana/sample01.csv",
				r: func(one []string, lineno int) {
					fmt.Println(lineno, one[3])
				},
				mode: ModeCsvSJIS,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadCSV(tt.args.fnm, tt.args.r, tt.args.mode)
		})
	}
}

func TestLoadCSVArr(t *testing.T) {
	type args struct {
		fnms []string
		r    CsvCallback
		mode int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadCSVArr(tt.args.fnms, tt.args.r, tt.args.mode)
		})
	}
}

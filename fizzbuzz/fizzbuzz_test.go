package fizzbuzz

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func fizzBuzz(n int) []string {
	return FizzBuzz(3, 5, n, "Fizz", "Buzz")
}

func Test_fizzBuzz(t *testing.T) {
	type args struct {
		n int
	}

	tests := map[string]struct {
		args args
		want []string
	}{
		"1":  {args{1}, []string{"1"}},
		"2":  {args{2}, []string{"1", "2"}},
		"3":  {args{3}, []string{"1", "2", "Fizz"}},
		"4":  {args{4}, []string{"1", "2", "Fizz", "4"}},
		"5":  {args{5}, []string{"1", "2", "Fizz", "4", "Buzz"}},
		"6":  {args{6}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz"}},
		"7":  {args{7}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7"}},
		"8":  {args{8}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8"}},
		"9":  {args{9}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz"}},
		"10": {args{10}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz"}},
		"11": {args{11}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11"}},
		"12": {args{12}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz"}},
		"13": {args{13}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13"}},
		"14": {args{14}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14"}},
		"15": {args{15}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"}},
		"16": {args{16}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz", "16"}},
		"17": {args{17}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz", "16", "17"}},
		"18": {args{18}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz", "16", "17", "Fizz"}},
		"19": {args{19}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz", "16", "17", "Fizz", "19"}},
		"20": {args{20}, []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz", "16", "17", "Fizz", "19", "Buzz"}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := fizzBuzz(tt.args.n)
			if !cmp.Equal(got, tt.want) {
				t.Error(cmp.Diff(got, tt.want))
			}
		})
	}
}

func BenchmarkFizzBuzz(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fizzBuzz(1_000_000)
	}
}

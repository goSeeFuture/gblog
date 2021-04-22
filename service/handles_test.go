package service

import (
	"reflect"
	"testing"
)

func Test_pages(t *testing.T) {
	const hellip = "&hellip;"

	type args struct {
		maxPage int
		curPage int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"1页", args{1, 1}, []string{"1"}},
		{"2页1", args{2, 1}, []string{"1", "2"}},
		{"2页2", args{2, 2}, []string{"1", "2"}},
		{"3页1", args{3, 1}, []string{"1", "2", "3"}},
		{"3页2", args{3, 2}, []string{"1", "2", "3"}},
		{"3页3", args{3, 3}, []string{"1", "2", "3"}},
		{"4页1", args{4, 1}, []string{"1", "2", hellip, "4"}},
		{"4页2", args{4, 2}, []string{"1", "2", "3", "4"}},
		{"4页3", args{4, 3}, []string{"1", "2", "3", "4"}},
		{"4页4", args{4, 4}, []string{"1", hellip, "3", "4"}},
		{"5页1", args{5, 1}, []string{"1", "2", hellip, "5"}},
		{"5页2", args{5, 2}, []string{"1", "2", "3", hellip, "5"}},
		{"5页3", args{5, 3}, []string{"1", "2", "3", "4", "5"}},
		{"5页4", args{5, 4}, []string{"1", hellip, "3", "4", "5"}},
		{"5页5", args{5, 5}, []string{"1", hellip, "4", "5"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pages(tt.args.maxPage, tt.args.curPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pages() = %v, want %v", got, tt.want)
			}
		})
	}
}
